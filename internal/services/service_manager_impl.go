package services

import (
	"example.com/api/config"
	"example.com/api/internal/repository"
	"example.com/api/internal/services/chat"
	"example.com/api/internal/services/hashing"
	"example.com/api/internal/storage"
	"example.com/api/internal/storage/cache"
	"example.com/api/pkg/logging"
)

type ServiceManager struct {
	repoManager  repository.IRepositoryManager
	logger       logging.ILogger
	config       config.Config
	user         IUserService
	chat         chat.IChatService
	auth         IAuthService
	hash         hashing.IHashService
	tokenStorage storage.ITokenStorage
	cacheStorage cache.ICacheService
}

func NewServiceManager(
	r repository.IRepositoryManager,
	l logging.ILogger,
	cfg config.Config,
) IServiceManager {

	return &ServiceManager{
		repoManager: r,
		logger:      l,
		config:      cfg,
	}
}

func (s *ServiceManager) User() IUserService {
	if s.user == nil {
		s.user = NewUserService(s.repoManager, s.logger, s.Hash())
	}
	return s.user
}

func (s *ServiceManager) Chat() chat.IChatService {
	if s.chat == nil {
		s.chat = chat.NewChatService(s.repoManager, s.logger)
	}
	return s.chat
}

func (s *ServiceManager) Auth() IAuthService {
	if s.auth == nil {
		s.auth = NewAuthService(
			s.config.JWT,
			s.Hash(),
			s.User(),
			s.logger,
			s.TokenStorage(),
		)
	}
	return s.auth
}

func (s *ServiceManager) Hash() hashing.IHashService {
	if s.hash == nil {
		s.hash = hashing.New()
	}
	return s.hash
}

func (s *ServiceManager) TokenStorage() storage.ITokenStorage {
	if s.tokenStorage == nil {
		tokenRedis := storage.NewRedisClient(&s.config.Redis, s.config.Redis.TokenStorage.DB)
		s.tokenStorage = storage.NewRedisTokenStorage(tokenRedis, s.config.Redis.TokenStorage.KeyPrefix)
	}
	return s.tokenStorage
}

func (s *ServiceManager) CacheStorage() cache.ICacheService {
	if s.cacheStorage == nil {
		cacheRedis := storage.NewRedisClient(&s.config.Redis, s.config.Redis.CacheStorage.DB)
		s.cacheStorage = cache.NewRedisCache(cacheRedis, s.config.Redis.CacheStorage.KeyPrefix)
	}
	return s.cacheStorage
}
