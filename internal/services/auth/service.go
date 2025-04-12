package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/models"
)

const (
	cacheExpTime = 10 * time.Minute

	userCacheKeyPrefix    = "auth:user:username"
	sessionCacheKeyPrefix = "auth:session:sessionID"
)

type userRepository interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type sessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, id int64) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
	MarkSessionAsRevoked(ctx context.Context, refreshToken string) error
}

type serv struct {
	userRepo    userRepository
	sessionRepo sessionRepository
	cache       cache.RedisClient
	jwtConfig   config.JWTConfig
}

func NewService(
	userRepo userRepository,
	sessionRepo sessionRepository,
	cache cache.RedisClient,
	jwtConfig config.JWTConfig) *serv {
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
