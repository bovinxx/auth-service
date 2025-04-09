package auth

import (
	"fmt"
	"time"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	sessionRepo "github.com/bovinxx/auth-service/internal/repository/session"
	userRepo "github.com/bovinxx/auth-service/internal/repository/user"
	authService "github.com/bovinxx/auth-service/internal/services/auth"
)

const (
	cacheExpTime = 10 * time.Minute

	userCacheKeyPrefix    = "auth:user:username"
	sessionCacheKeyPrefix = "auth:session:sessionID"
)

type serv struct {
	userRepo    userRepo.Repository
	sessionRepo sessionRepo.Repository
	cache       cache.RedisClient
	jwtConfig   config.JWTConfig
}

func NewService(
	userRepo userRepo.Repository,
	sessionRepo sessionRepo.Repository,
	cache cache.RedisClient,
	jwtConfig config.JWTConfig) authService.Service {
	return &serv{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		cache:       cache,
		jwtConfig:   jwtConfig,
	}
}

func newUserCacheKey(prefix, username string) string {
	return fmt.Sprintf("%s:%s", prefix, username)
}
