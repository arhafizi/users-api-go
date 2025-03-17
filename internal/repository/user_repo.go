package repository

import (
	dbCtx "example.com/api/internal/repository/db"
)

type IUserRepo interface {
	GetByID(ctx Ctx, id int32) (User, error)

	GetByUsername(ctx Ctx, username string) (User, error)

	GetByEmail(ctx Ctx, email string) (User, error)

	GetAll(ctx Ctx, arg dbCtx.ListUsersParams) ([]User, error)

	Create(ctx Ctx, arg dbCtx.CreateUserParams) (User, error)

	UpdateFull(ctx Ctx, arg dbCtx.UpdateUserFullParams) (User, error)

	UpdatePartial(ctx Ctx, arg dbCtx.UpdateUserPartialParams) (User, error)

	SoftDelete(ctx Ctx, id int32) error
}
