package repository

import (
	"context"

	dbCtx "example.com/api/internal/repository/db"
)

type IUserRepo interface {
	GetByID(ctx context.Context, id int32) (dbCtx.User, error)
	GetByUsername(ctx context.Context, username string) (dbCtx.User, error)
	GetAll(ctx context.Context, arg dbCtx.ListUsersParams) ([]dbCtx.User, error)
	Create(ctx context.Context, arg dbCtx.CreateUserParams) (dbCtx.User, error)
	UpdateFull(ctx context.Context, arg dbCtx.UpdateUserFullParams) (dbCtx.User, error)
	UpdatePartial(ctx context.Context, arg dbCtx.UpdateUserPartialParams) (dbCtx.User, error)
	SoftDelete(ctx context.Context, id int32) error
}
