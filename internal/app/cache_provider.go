package app

import (
	"context"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/client/cache/redis"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/logger"
	redigo "github.com/gomodule/redigo/redis"
)

type cacheProvider struct {
	redisConfig config.RedisConfig
	redisClient cache.RedisClient
	redisPool   *redigo.Pool
}

func newCacheProvider() *cacheProvider {
	return &cacheProvider{}
}

func (s *cacheProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			logger.Fatal("failed to create a new redis config", logger.Err(err))
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *cacheProvider) RedisClient(ctx context.Context) cache.RedisClient {
	if s.redisClient == nil {
		client := redis.NewClient(s.RedisPool(), s.RedisConfig())

		err := client.Ping(ctx)
		if err != nil {
			logger.Fatal("failed to ping redis client", logger.Err(err))
		}

		s.redisClient = client
	}

	return s.redisClient
}

func (s *cacheProvider) RedisPool() *redigo.Pool {
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
