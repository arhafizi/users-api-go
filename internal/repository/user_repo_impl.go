package repository

import (
	"context"
	"database/sql"

	dbCtx "example.com/api/internal/repository/db"
)

type UserRepo struct {
	q *dbCtx.Queries
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &UserRepo{
		q: dbCtx.New(db),
	}
}

func (u *UserRepo) Create(ctx context.Context, arg dbCtx.CreateUserParams) (dbCtx.User, error) {
	return u.q.CreateUser(ctx, arg)
}

func (u *UserRepo) GetAll(ctx context.Context, arg dbCtx.ListUsersParams) ([]dbCtx.User, error) {
	return u.q.ListUsers(ctx, arg)
}

func (u *UserRepo) GetByID(ctx context.Context, id int32) (dbCtx.User, error) {
	return u.q.GetUserByID(ctx, id)
}

func (u *UserRepo) GetByUsername(ctx context.Context, username string) (dbCtx.User, error) {
	return u.q.GetUserByUsername(ctx, username)
}

func (u *UserRepo) SoftDelete(ctx context.Context, id int32) error {
	return u.q.SoftDeleteUser(ctx, id)
}

func (u *UserRepo) UpdateFull(ctx context.Context, arg dbCtx.UpdateUserFullParams) (dbCtx.User, error) {
	return u.q.UpdateUserFull(ctx, arg)
}

func (u *UserRepo) UpdatePartial(ctx context.Context, arg dbCtx.UpdateUserPartialParams) (dbCtx.User, error) {
	return u.q.UpdateUserPartial(ctx, arg)
}
