package services

import (
	"context"
	"errors"
	"log"

	"example.com/api/internal/repository"
	dbCtx "example.com/api/internal/repository/db"
)

type UserService struct {
	repo repository.IUserRepo
}

func NewUserService(repo repository.IUserRepo) IUserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, arg dbCtx.CreateUserParams) (*dbCtx.User, error) {
	user, err := s.repo.Create(ctx, arg)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, errors.New("failed to create user")
	}
	return &user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int32) (*dbCtx.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Failed to fetch user by ID: %v", err)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*dbCtx.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("Failed to fetch user by username: %v", err)
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserService) GetAll(ctx context.Context, arg dbCtx.ListUsersParams) ([]dbCtx.User, error) {
	users, err := s.repo.GetAll(ctx, arg)
	if err != nil {
		log.Printf("Failed to fetch all users: %v", err)
		return nil, errors.New("failed to fetch users")
	}
	return users, nil
}

func (s *UserService) SoftDelete(ctx context.Context, id int32) error {
	err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		log.Printf("Failed to soft delete user: %v", err)
		return errors.New("failed to delete user")
	}
	return nil
}

func (s *UserService) UpdateFull(ctx context.Context, arg dbCtx.UpdateUserFullParams) (*dbCtx.User, error) {
	user, err := s.repo.UpdateFull(ctx, arg)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, errors.New("failed to update user")
	}
	return &user, nil
}

func (s *UserService) UpdatePartial(ctx context.Context, arg dbCtx.UpdateUserPartialParams) (*dbCtx.User, error) {
	user, err := s.repo.UpdatePartial(ctx, arg)
	if err != nil {
		log.Printf("Failed to partially update user: %v", err)
		return nil, errors.New("failed to update user")
	}
	return &user, nil
}
