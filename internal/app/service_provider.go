package app

import (
	"context"

	accessAPI "github.com/bovinxx/auth-service/internal/api/access"
	authAPI "github.com/bovinxx/auth-service/internal/api/auth"
	userAPI "github.com/bovinxx/auth-service/internal/api/user"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/logger"
	"github.com/bovinxx/auth-service/internal/models"
	sessionRepoImpl "github.com/bovinxx/auth-service/internal/repository/session"
	userRepoImpl "github.com/bovinxx/auth-service/internal/repository/user"
	accessServiceImpl "github.com/bovinxx/auth-service/internal/services/access"
	authServiceImpl "github.com/bovinxx/auth-service/internal/services/auth"
	userServiceImpl "github.com/bovinxx/auth-service/internal/services/user"
	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userService interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}

type authService interface {
	Login(ctx context.Context, req *models.User) (string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetRefreshToken(ctx context.Context, token *models.Token) (*models.Token, error)
	GetAccessToken(ctx context.Context, token *models.Token) (*models.Token, error)
}

type accessService interface {
	Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error)
}

type userRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}

type sessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, id int64) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
	GetAllUserSessions(ctx context.Context, username string) ([]*models.Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
	MarkSessionAsRevoked(ctx context.Context, refreshToken string) error
}

type serviceProvider struct {
	dbProvider    *dbProvider
	cacheProvider *cacheProvider

	grpcConfig       config.GRPCConfig
	httpConfig       config.HTTPConfig
	prometheusConfig config.PrometheusConfig
	jwtConfig        config.JWTConfig

	userRepo    userRepository
	sessionRepo sessionRepository

	userService   userService
	authService   authService
	accessService accessService

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

func (s *serviceProvider) UserRepository(ctx context.Context) userRepository {
	if s.userRepo == nil {
		repo, err := userRepoImpl.NewRepository(s.DBProvider().DBClient(ctx))
		if err != nil {
			logger.Fatal("failed to create user repository", logger.Err(err))
		}

		s.userRepo = repo
	}

	return s.userRepo
}

func (s *serviceProvider) SessionRepository(ctx context.Context) sessionRepository {
	if s.sessionRepo == nil {
		repo, err := sessionRepoImpl.NewRepository(s.DBProvider().DBClient(ctx))
		if err != nil {
			logger.Fatal("failed to create session repository", logger.Err(err))
		}

		s.sessionRepo = repo
	}

	return s.sessionRepo
}

func (s *serviceProvider) UserService(ctx context.Context) userService {
	if s.userService == nil {
		service := userServiceImpl.NewService(s.UserRepository(ctx), s.dbProvider.TxManager(ctx))
		s.userService = service
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) authService {
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

func (s *serviceProvider) AccessService(ctx context.Context) accessService {
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
