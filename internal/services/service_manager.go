package services

import (
	"example.com/api/internal/services/chat"
	"example.com/api/internal/services/hashing"
	"example.com/api/internal/storage"
	"example.com/api/internal/storage/cache"
)

type IServiceManager interface {
	User() IUserService
	Chat() chat.IChatService
	Auth() IAuthService
	Hash() hashing.IHashService
	TokenStorage() storage.ITokenStorage
	CacheStorage() cache.ICacheService
}
