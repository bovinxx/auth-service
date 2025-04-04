package app

import (
	"context"
	"log"

	userAPI "github.com/bovinxx/auth-service/internal/api/user"
	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/client/cache/redis"
	"github.com/bovinxx/auth-service/internal/client/db"
	"github.com/bovinxx/auth-service/internal/client/db/pg"
	"github.com/bovinxx/auth-service/internal/client/db/transaction"
	"github.com/bovinxx/auth-service/internal/closer"
	"github.com/bovinxx/auth-service/internal/config"
	repository "github.com/bovinxx/auth-service/internal/repository/user"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user/user"
	service "github.com/bovinxx/auth-service/internal/services"
	userService "github.com/bovinxx/auth-service/internal/services/user/user"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig         config.PGConfig
	grpcConfig       config.GRPCConfig
	httpConfig       config.HTTPConfig
	prometheusConfig config.PrometheusConfig
	redisConfig      config.RedisConfig

	dbClient    db.Client
	redisClient cache.RedisClient

	redisPool *redigo.Pool
	txManager db.TxManager

	userRepo repository.UserRepository

	userService service.UserService
	// authService   service.AuthService
	// accessService service.AccessService

	userImpl *userAPI.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		db, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}
		closer.Add(db.Close)
		s.dbClient = db
	}

	return s.dbClient
}

func (s *serviceProvider) RedisClient(ctx context.Context) cache.RedisClient {
	if s.redisClient == nil {
		client := redis.NewClient(s.RedisPool(), s.RedisConfig())

		err := client.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to create a new redis client: %v", err)
		}

		s.redisClient = client
	}

	return s.redisClient
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		pool := &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}

		s.redisPool = pool
	}

	return s.redisPool
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to create a new postgres config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.DBClient(ctx)
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to create a new grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to create a new http config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		cfg, err := config.NewPrometheusConfig()
		if err != nil {
			log.Fatalf("failed to create a new prometheus config: %v", err)
		}

		s.prometheusConfig = cfg
	}

	return s.prometheusConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to create a new redis config: %v", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		repo, err := userRepo.NewRepository(s.DBClient(ctx))
		if err != nil {
			log.Fatalf("failed to create repository: %v", err)
		}

		s.userRepo = repo
	}

	return s.userRepo
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		service := userService.NewService(s.UserRepository(ctx), s.TxManager(ctx))

		s.userService = service
	}

	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *userAPI.Implementation {
	if s.userImpl == nil {
		impl := userAPI.NewImplementation(s.UserService(ctx))

		s.userImpl = impl
	}

	return s.userImpl
}
