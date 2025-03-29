package services

import (
	"context"

	dto "example.com/api/internal/contracts"
	dbCtx "example.com/api/internal/repository/db"
)

type IUserService interface {
	GetByID(ctx context.Context, id int32) (*dbCtx.User, error)

	Create(ctx context.Context, arg dto.CreateUserReq) (*dto.UserResponse, error)

	GetByUsername(ctx context.Context, username string) (*dbCtx.User, error)

	GetByEmail(ctx context.Context, email string) (*dbCtx.User, error)

	GetAll(ctx context.Context, arg dto.ListUsersParams) ([]dbCtx.User, error)

	SoftDelete(ctx context.Context, id int32) error

	UpdateFull(ctx context.Context, arg dto.UpdateUserFullReq) (*dbCtx.User, error)

	UpdatePartial(ctx context.Context, arg dto.UpdateUserPartialReq) (*dbCtx.User, error)
}
