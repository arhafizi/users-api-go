package repository

import (
	"context"
	"database/sql"
	"errors"

	dbCtx "example.com/api/internal/repository/db"
)

type RepositoryManager struct {
	db       dbCtx.DBTX
	userRepo IUserRepo
	chatRepo IChatRepo
}

func NewRepositoryManager(db dbCtx.DBTX) IRepositoryManager {
	return &RepositoryManager{
		db: db,
	}
}

func (rm *RepositoryManager) WithTx(ctx context.Context, fn func(IRepositoryManager) error) error {
	db, ok := rm.db.(*sql.DB)
	if !ok {
		return errors.New("WithTx requires a *sql.DB as the base connection")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txRepoMgr := NewRepositoryManager(tx)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	if err := fn(txRepoMgr); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *RepositoryManager) User() IUserRepo {
	if r.userRepo == nil {
		r.userRepo = NewUserRepo(r.db)
	}
	return r.userRepo
}

func (r *RepositoryManager) Chat() IChatRepo {
	if r.chatRepo == nil {
		r.chatRepo = NewChatRepo(r.db)
	}
	return r.chatRepo
}
