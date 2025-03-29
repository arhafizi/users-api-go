package services

import (
	"context"

	"example.com/api/config"
	"example.com/api/internal/repository"
	"example.com/api/internal/services/hashing"
	"example.com/api/pkg/logging"
)

type ServiceManager struct {
	repoManager repository.IRepositoryManager
	logger      logging.ILogger
	jwtConfig   config.JWTConfig
	userSvc     IUserService
	authSvc     IAuthService
	hashSvc     hashing.IHashService
}

func NewServiceManager(r repository.IRepositoryManager, l logging.ILogger, jwt config.JWTConfig) IServiceManager {
	return &ServiceManager{
		repoManager: r,
		logger:      l,
		jwtConfig:   jwt,
	}
}

func (s *ServiceManager) User() IUserService {
	if s.userSvc == nil {
		s.userSvc = NewUserService(s.repoManager, s.logger, s.Hash())
	}
	return s.userSvc
}

func (s *ServiceManager) Auth() IAuthService {
	if s.authSvc == nil {
		s.authSvc = NewAuthService(s.jwtConfig, s.Hash(), s.User(), s.logger)
	}
	return s.authSvc
}

func (s *ServiceManager) Hash() hashing.IHashService {
	if s.hashSvc == nil {
		s.hashSvc = hashing.New()
	}
	return s.hashSvc
}

func (s *ServiceManager) WithTransaction(ctx context.Context, fn func(txService IServiceManager) error) (err error) {
	txRepo, err := s.repoManager.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Create a transactional ServiceManager instance that uses the transactional RepositoryManager.
	txService := NewServiceManager(txRepo, s.logger, s.jwtConfig)

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

// func (s *UserService) CreateUserWithProfile(ctx context.Context, userParams dbCtx.CreateUserParams, profileParams dbCtx.CreateProfileParams) error {
//     return s.serviceManager.WithTransaction(ctx, func(txService IServiceManager) error {
//         // Create user using the transactional context.
//         _, err := txService.User().Create(ctx, userParams)
//         if err != nil {
//             return err
//         }

//         // Example: Create profile using the same transactional context.
//         // (Assuming txService.Profile() is implemented similarly.)
//         if err := txService.Profile().Create(ctx, profileParams); err != nil {
//             return err
//         }

//         // All operations share the same transaction: commit if nil is returned.
//         return nil
//     })
// }
