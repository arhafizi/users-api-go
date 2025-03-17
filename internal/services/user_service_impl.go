package services

import (
	"context"
	"database/sql"
	"errors"

	"example.com/api/internal/repository"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/pkg/logging"
)

type UserService struct {
	repo   repository.IRepositoryManager
	logger logging.ILogger
}

func NewUserService(r repository.IRepositoryManager, l logging.ILogger) *UserService {
	return &UserService{
		repo:   r,
		logger: l,
	}
}

func (s *UserService) Create(ctx context.Context, arg dbCtx.CreateUserParams) (*dbCtx.User, error) {
	user, err := s.repo.User().Create(ctx, arg)
	if err != nil {
		s.logger.Error(
			logging.Internal, logging.FailedToCreateUser, "Failed to create user",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
			},
		)
		return nil, errors.New("failed to create user")
	}
	return &user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int32) (*dbCtx.User, error) {
	user, err := s.repo.User().GetByID(ctx, id)
	if err != nil {
		s.logger.Error(
			logging.Postgres, logging.Select, "Failed to fetch user by ID",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"userID":             id,
			},
		)
		return nil, errors.New("user not found")
	}
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

func (s *UserService) GetAll(ctx context.Context, arg dbCtx.ListUsersParams) ([]dbCtx.User, error) {
	users, err := s.repo.User().GetAll(ctx, arg)
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
	err := s.repo.User().SoftDelete(ctx, id)
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
	return nil
}

func (s *UserService) UpdateFull(ctx context.Context, arg dbCtx.UpdateUserFullParams) (*dbCtx.User, error) {
	user, err := s.repo.User().UpdateFull(ctx, arg)
	if err != nil {
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
		return nil, errors.New("failed to update user")
	}
	return &user, nil
}

func (s *UserService) UpdatePartial(ctx context.Context, arg dbCtx.UpdateUserPartialParams) (*dbCtx.User, error) {
	user, err := s.repo.User().UpdatePartial(ctx, arg)
	if err != nil {
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
			logging.Postgres, logging.Update, "Failed to partially update user",
			map[logging.ExtraKey]any{
				logging.ErrorMessage: err.Error(),
				"userID":             arg.ID,
			},
		)
		return nil, errors.New("failed to update user")
	}
	return &user, nil
}
