package app

import (
	"context"

	accessAPI "github.com/bovinxx/auth-service/internal/api/access"
	authAPI "github.com/bovinxx/auth-service/internal/api/auth"
	userAPI "github.com/bovinxx/auth-service/internal/api/user"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/logger"
	sessionRepo "github.com/bovinxx/auth-service/internal/repository/session"
	sessionRepoImpl "github.com/bovinxx/auth-service/internal/repository/session/session"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	userRepoImpl "github.com/bovinxx/auth-service/internal/repository/user/user"
	accessService "github.com/bovinxx/auth-service/internal/services/access"
	accessServiceImpl "github.com/bovinxx/auth-service/internal/services/access/access"
	authService "github.com/bovinxx/auth-service/internal/services/auth"
	authServiceImpl "github.com/bovinxx/auth-service/internal/services/auth/auth"
	userService "github.com/bovinxx/auth-service/internal/services/user"
	userServiceImpl "github.com/bovinxx/auth-service/internal/services/user/user"
)

type serviceProvider struct {
	dbProvider    *dbProvider
	cacheProvider *cacheProvider

	grpcConfig       config.GRPCConfig
	httpConfig       config.HTTPConfig
	prometheusConfig config.PrometheusConfig
	jwtConfig        config.JWTConfig

	userRepo    userRepo.Repository
	sessionRepo sessionRepo.Repository

	userService   userService.Service
	authService   authService.Service
	accessService accessService.Service

	userImpl   *userAPI.Implementation
	authImpl   *authAPI.Implementation
	accessImpl *accessAPI.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DBProvider() *dbProvider {
	if s.dbProvider == nil {
		provider := newDBProvider()
		s.dbProvider = provider
	}

	return s.dbProvider
}

func (s *serviceProvider) CacheProvider() *cacheProvider {
	if s.cacheProvider == nil {
		provider := newCacheProvider()
		s.cacheProvider = provider
	}

	return s.cacheProvider
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			logger.Fatal("failed to create a new grpc config", logger.Err(err))
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			logger.Fatal("failed to create a new http config", logger.Err(err))
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		cfg, err := config.NewPrometheusConfig()
		if err != nil {
			logger.Fatal("failed to create a new prometheus config", logger.Err(err))
		}

		s.prometheusConfig = cfg
	}

	return s.prometheusConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			logger.Fatal("failed to create jwt config", logger.Err(err))
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) UserRepository(ctx context.Context) userRepo.Repository {
	if s.userRepo == nil {
		repo, err := userRepoImpl.NewRepository(s.DBProvider().DBClient(ctx))
		if err != nil {
			logger.Fatal("failed to create user repository", logger.Err(err))
		}

		s.userRepo = repo
	}

	return s.userRepo
}

func (s *serviceProvider) SessionRepository(ctx context.Context) sessionRepo.Repository {
	if s.sessionRepo == nil {
		repo, err := sessionRepoImpl.NewRepository(s.DBProvider().DBClient(ctx))
		if err != nil {
			logger.Fatal("failed to create session repository", logger.Err(err))
		}

		s.sessionRepo = repo
	}

	return s.sessionRepo
}

func (s *serviceProvider) UserService(ctx context.Context) userService.Service {
	if s.userService == nil {
		service := userServiceImpl.NewService(s.UserRepository(ctx), s.dbProvider.TxManager(ctx))
		s.userService = service
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) authService.Service {
	if s.authService == nil {
		service := authServiceImpl.NewService(
			s.UserRepository(ctx),
			s.SessionRepository(ctx),
			s.CacheProvider().RedisClient(ctx),
			s.JWTConfig(),
		)
		s.authService = service
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) accessService.Service {
	if s.accessService == nil {
		service := accessServiceImpl.NewService(
			s.UserRepository(ctx),
			s.SessionRepository(ctx),
			s.CacheProvider().RedisClient(ctx),
			s.JWTConfig(),
		)
		s.accessService = service
	}

	return s.accessService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *userAPI.Implementation {
	if s.userImpl == nil {
		impl := userAPI.NewImplementation(s.UserService(ctx))
		s.userImpl = impl
	}

	return s.userImpl
}

func (s *serviceProvider) AuthImplementation(ctx context.Context) *authAPI.Implementation {
	if s.authImpl == nil {
		impl := authAPI.NewImplementation(s.AuthService(ctx))
		s.authImpl = impl
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImplementation(ctx context.Context) *accessAPI.Implementation {
	if s.accessImpl == nil {
		impl := accessAPI.NewImplementation(s.AccessService(ctx))
		s.accessImpl = impl
	}

	return s.accessImpl
}
