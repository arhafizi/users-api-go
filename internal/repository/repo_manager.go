package repository

import "context"

type IRepositoryManager interface {
	User() IUserRepo
	Chat() IChatRepo

	BeginTx(ctx context.Context) (IRepositoryManager, error)
	Commit() error
	Rollback() error
}
