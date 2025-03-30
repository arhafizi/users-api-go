package services

import (
	"context"

	"example.com/api/internal/services/hashing"
	"example.com/api/internal/storage"
	"example.com/api/internal/storage/cache"
)

type IServiceManager interface {
	WithTransaction(ctx context.Context, fn func(txService IServiceManager) error) (err error)
	User() IUserService
	Auth() IAuthService
	Hash() hashing.IHashService
	TokenStorage() storage.ITokenStorage
	CacheStorage() cache.ICacheService
}
