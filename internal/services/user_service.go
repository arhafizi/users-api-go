package services

import (
	"context"

	dbCtx "example.com/api/internal/repository/db"
)

type IUserService interface {
	GetByID(ctx context.Context, id int32) (*dbCtx.User, error)

	Create(ctx context.Context, arg dbCtx.CreateUserParams) (*dbCtx.User, error)

	GetByUsername(ctx context.Context, username string) (*dbCtx.User, error)

	GetByEmail(ctx context.Context, email string) (*dbCtx.User, error)

	GetAll(ctx context.Context, arg dbCtx.ListUsersParams) ([]dbCtx.User, error)

	SoftDelete(ctx context.Context, id int32) error

	UpdateFull(ctx context.Context, arg dbCtx.UpdateUserFullParams) (*dbCtx.User, error)

	UpdatePartial(ctx context.Context, arg dbCtx.UpdateUserPartialParams) (*dbCtx.User, error)
}
