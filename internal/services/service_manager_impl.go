package services

import (
	"context"

	"example.com/api/config"
	"example.com/api/internal/repository"
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

func (s *ServiceManager) WithTransaction(ctx context.Context, fn func(txService IServiceManager) error) (err error) {
	txRepo, err := s.repoManager.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Create a transactional ServiceManager instance that uses the transactional RepositoryManager.
	txService := NewServiceManager(txRepo, s.logger, s.config)

	// Defer commit/rollback logic.
	defer func() {
		if p := recover(); p != nil {
			// Roll back on panic and re-panic.
			_ = txRepo.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := txRepo.Rollback(); rbErr != nil {
				s.logger.Errorf("Transaction rollback error: %v", rbErr)
			}
		} else {
			err = txRepo.Commit()
		}
	}()

	// Execute the provided transactional function.
	err = fn(txService)
	return err
}
