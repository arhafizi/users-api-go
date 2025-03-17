package services

import (
	"context"

	"example.com/api/internal/services/hashing"
)

type IServiceManager interface {
	WithTransaction(ctx context.Context, fn func(txService IServiceManager) error) (err error)
	User() IUserService
	Auth() IAuthService
	Hash() hashing.IHashService
}
