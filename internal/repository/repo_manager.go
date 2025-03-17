package repository

import "context"

type IRepositoryManager interface {
	User() IUserRepo

	BeginTx(ctx context.Context) (IRepositoryManager, error)
	Commit() error
	Rollback() error
}
