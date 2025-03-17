package repository

import (
	dbCtx "example.com/api/internal/repository/db"
)

type UserRepo struct {
	q *dbCtx.Queries
}

func NewUserRepo(db dbCtx.DBTX) IUserRepo {
	return &UserRepo{
		q: dbCtx.New(db),
	}
}
func (u *UserRepo) Create(ctx Ctx, arg dbCtx.CreateUserParams) (User, error) {
	return u.q.CreateUser(ctx, arg)
}

func (u *UserRepo) GetAll(ctx Ctx, arg dbCtx.ListUsersParams) ([]User, error) {
	return u.q.ListUsers(ctx, arg)
}

func (u *UserRepo) GetByID(ctx Ctx, id int32) (User, error) {
	return u.q.GetUserByID(ctx, id)
}

func (u *UserRepo) GetByUsername(ctx Ctx, username string) (User, error) {
	return u.q.GetUserByUsername(ctx, username)
}

func (u *UserRepo) GetByEmail(ctx Ctx, username string) (User, error) {
	return u.q.GetUserByEmail(ctx, username)
}

func (u *UserRepo) SoftDelete(ctx Ctx, id int32) error {
	return u.q.SoftDeleteUser(ctx, id)
}

func (u *UserRepo) UpdateFull(ctx Ctx, arg dbCtx.UpdateUserFullParams) (User, error) {
	return u.q.UpdateUserFull(ctx, arg)
}

func (u *UserRepo) UpdatePartial(ctx Ctx, arg dbCtx.UpdateUserPartialParams) (User, error) {
	return u.q.UpdateUserPartial(ctx, arg)
}
