package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *RepositoryManager) BeginTx(ctx context.Context) (IRepositoryManager, error) {
	if _, ok := r.db.(*sql.Tx); ok {
		return nil, errors.New("already in a transaction; nested transactions are not supported")
	}

	db, ok := r.db.(*sql.DB)
	if !ok {
		return nil, errors.New("BeginTx: underlying db is not *sql.DB")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("BeginTx: %w", err)
	}

	return NewRepositoryManager(tx), nil
}

func (r *RepositoryManager) Commit() error {
	tx, ok := r.db.(*sql.Tx)
	if !ok {
		return errors.New("Commit: not operating in a transaction")
	}
	return tx.Commit()
}

func (r *RepositoryManager) Rollback() error {
	tx, ok := r.db.(*sql.Tx)
	if !ok {
		return errors.New("Rollback: not operating in a transaction")
	}
	return tx.Rollback()
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
