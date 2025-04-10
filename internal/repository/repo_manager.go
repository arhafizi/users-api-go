package repository

import (
	"context"
)

type IRepositoryManager interface {
	User() IUserRepo
	Chat() IChatRepo
	WithTx(context.Context, func(IRepositoryManager) error) error
}
