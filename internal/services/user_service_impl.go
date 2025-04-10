package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	dto "example.com/api/internal/contracts"
	contracts "example.com/api/internal/contracts/errors"
	"example.com/api/internal/repository"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services/hashing"
	"example.com/api/pkg/logging"
	"example.com/api/pkg/metrics"
	"github.com/lib/pq"
)

type UserService struct {
	hashService hashing.IHashService
	repo        repository.IRepositoryManager
	logger      logging.ILogger
}

func NewUserService(r repository.IRepositoryManager, l logging.ILogger, h hashing.IHashService) *UserService {
	return &UserService{
		repo:        r,
		logger:      l,
		hashService: h,
	}
}

func (s *UserService) Create(ctx context.Context, arg dto.CreateUserReq) (*dto.UserResponse, error) {
	metrics.DbCall.WithLabelValues("User", "Create", "started").Inc()
	hashedPassword, err := s.hashService.Hash(arg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	arg.Password = hashedPassword
	params := mapCreateUserReqToParams(arg)
	user, err := s.repo.User().Create(ctx, params)
	if err != nil {
		metrics.DbCall.WithLabelValues("User", "Create", "error").Inc()

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key":
				s.logger.Warn(
					logging.Validation, logging.FailedToCreateUser, "Username already exists",
					map[logging.ExtraKey]any{logging.RequestBody: arg.Username},
				)
				return nil, &contracts.UsernameExistsError{Username: arg.Username}

			case "users_email_key":
				s.logger.Warn(
					logging.Validation, logging.FailedToCreateUser, "Email already exists",
					map[logging.ExtraKey]any{logging.RequestBody: arg.Email},
				)
				return nil, &contracts.EmailExistsError{Email: arg.Email}
			}
		}

		s.logger.Error(
			logging.Internal, logging.FailedToCreateUser, "Failed to create user",
			map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()},
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	metrics.DbCall.WithLabelValues("User", "Create", "success").Inc()
	userResponse := mapUserToResponse(user)
	return &userResponse, nil
}

func (s *UserService) GetByID(ctx context.Context, id int32) (*dbCtx.User, error) {
	metrics.DbCall.WithLabelValues("User", "GetById", "started").Inc()

	user, err := s.repo.User().GetByID(ctx, id)
	if err != nil {
		metrics.DbCall.WithLabelValues("User", "Get", "error").Inc()

		s.logger.Error(
			logging.Postgres, logging.Select, "Failed to fetch user by ID",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"userID":             id,
			},
		)
		return nil, errors.New("user not found")
	}
	metrics.DbCall.WithLabelValues("User", "Get", "success").Inc()

	return &user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*dbCtx.User, error) {
	user, err := s.repo.User().GetByUsername(ctx, username)
	if err != nil {
		s.logger.Error(
			logging.Postgres, logging.Select, "Failed to fetch user by username",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"username":           username,
			},
		)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*dbCtx.User, error) {
	user, err := s.repo.User().GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error(
			logging.Postgres, logging.Select, "Failed to fetch user by email",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"email":              email,
			},
		)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserService) GetAll(ctx context.Context, arg dto.ListUsersParams) ([]dbCtx.User, error) {
	params := mapListUsersReqToParams(arg)

	users, err := s.repo.User().GetAll(ctx, params)
	if err != nil {
		s.logger.Error(
			logging.Postgres, logging.Select, "Failed to fetch all users",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
			},
		)
		return nil, errors.New("failed to fetch users")
	}
	return users, nil
}

func (s *UserService) SoftDelete(ctx context.Context, id int32) error {
	rowsAffected, err := s.repo.User().SoftDelete(ctx, id)
	if err != nil {
		s.logger.Error(
			logging.Postgres, logging.Delete, "Failed to soft delete user",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"userID":             id,
			},
		)
		return errors.New("failed to delete user")
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *UserService) UpdateFull(ctx context.Context, arg dto.UpdateUserFullReq) (*dbCtx.User, error) {
	hashedPassword, _ := s.hashService.Hash(arg.Password)
	arg.Password = hashedPassword
	params := mapUpdateUserFullReqToParams(arg)

	user, err := s.repo.User().UpdateFull(ctx, params)
	if err != nil {
		metrics.DbCall.WithLabelValues("User", "UpdateFull", "error").Inc()

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key":
				s.logger.Warn(
					logging.Validation, logging.Update, "Username already exists",
					map[logging.ExtraKey]any{logging.RequestBody: arg.Username},
				)
				return nil, &contracts.UsernameExistsError{Username: arg.Username}

			case "users_email_key":
				s.logger.Warn(
					logging.Validation, logging.Update, "Email already exists",
					map[logging.ExtraKey]any{logging.RequestBody: arg.Email},
				)
				return nil, &contracts.EmailExistsError{Email: arg.Email}
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error(
				logging.Postgres, logging.Update, "User not found",
				map[logging.ExtraKey]any{
					logging.ErrorMessage: err.Error(),
					"userID":             arg.ID,
				},
			)
			return nil, errors.New("user not found")
		}

		s.logger.Error(
			logging.Postgres, logging.Update, "Failed to update user",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"userID":             arg.ID,
			},
		)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	metrics.DbCall.WithLabelValues("User", "UpdateFull", "success").Inc()
	return &user, nil
}

func (s *UserService) UpdatePartial(ctx context.Context, arg dto.UpdateUserPartialReq) (*dbCtx.User, error) {
	if arg.Password != nil && *arg.Password != "" {
		hashedPassword, _ := s.hashService.Hash(*arg.Password)
		arg.Password = &hashedPassword
	}
	params := mapUpdateUserPartialReqToParams(arg)

	user, err := s.repo.User().UpdatePartial(ctx, params)
	if err != nil {
		metrics.DbCall.WithLabelValues("User", "UpdatePartial", "error").Inc()

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key":
				s.logger.Warn(
					logging.Validation, logging.Update, "Username already exists",
					map[logging.ExtraKey]any{logging.RequestBody: arg.Username},
				)
				return nil, &contracts.UsernameExistsError{Username: *arg.Username}

			case "users_email_key":
				s.logger.Warn(
					logging.Validation, logging.Update, "Email already exists",
					map[logging.ExtraKey]any{logging.RequestBody: arg.Email},
				)
				return nil, &contracts.EmailExistsError{Email: *arg.Email}
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error(
				logging.Postgres, logging.Update, "User not found",
				map[logging.ExtraKey]any{
					logging.ErrorMessage: err.Error(),
					"userID":             arg.ID,
				},
			)
			return nil, errors.New("user not found")
		}

		s.logger.Error(
			logging.Postgres, logging.Update, "Failed to update user",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"userID":             arg.ID,
			},
		)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	metrics.DbCall.WithLabelValues("User", "UpdatePartial", "success").Inc()
	return &user, nil
}

func mapCreateUserReqToParams(dto dto.CreateUserReq) dbCtx.CreateUserParams {
	return dbCtx.CreateUserParams{
		Username:     dto.Username,
		Email:        dto.Email,
		FullName:     dto.FullName,
		PasswordHash: dto.Password,
	}
}

func mapUpdateUserFullReqToParams(dto dto.UpdateUserFullReq) dbCtx.UpdateUserFullParams {
	return dbCtx.UpdateUserFullParams{
		ID:           dto.ID,
		Username:     dto.Username,
		Email:        dto.Email,
		FullName:     dto.FullName,
		PasswordHash: dto.Password,
	}
}

func mapUpdateUserPartialReqToParams(dto dto.UpdateUserPartialReq) dbCtx.UpdateUserPartialParams {
	params := dbCtx.UpdateUserPartialParams{
		ID: dto.ID,
	}
	if dto.Username != nil {
		params.Username = dto.Username
	}
	if dto.Email != nil {
		params.Email = dto.Email
	}
	if dto.FullName != nil {
		params.FullName = dto.FullName
	}
	if dto.Password != nil {
		params.PasswordHash = dto.Password
	}
	return params
}

func mapListUsersReqToParams(dto dto.ListUsersParams) dbCtx.ListUsersParams {
	return dbCtx.ListUsersParams{
		Limit:  dto.Limit,
		Offset: dto.Offset,
	}
}

func mapUserToResponse(user dbCtx.User) dto.UserResponse {
	var createdAt *time.Time

	if user.CreatedAt.Valid {
		createdAt = &user.CreatedAt.Time
	}
	// if user.UpdatedAt.Valid {
	// 	updatedAt = &user.UpdatedAt.Time
	// }
	// if user.DeletedAt.Valid {
	// 	deletedAt = &user.DeletedAt.Time
	// }

	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: createdAt,
		// UpdatedAt: updatedAt,
		// DeletedAt: deletedAt,
	}
}

// example of a transaction
func (s *UserService) UpdateUserTx(ctx context.Context, args dbCtx.UpdateUserFullParams) (*dbCtx.User, error) {
	var updatedUser dbCtx.User
	err := s.repo.WithTx(ctx, func(txRM repository.IRepositoryManager) error {
		user, err := txRM.User().UpdateFull(ctx, args)
		if err != nil {
			return err
		}
		updatedUser = user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}
