package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/bovinxx/auth-service/internal/models"
	repoerrs "github.com/bovinxx/auth-service/internal/repository/session/errors"
	serverrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/pkg/errors"
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

func (s *serv) getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	cacheKey := newUserCacheKey(userCacheKeyPrefix, username)

	err := s.cache.GetStruct(ctx, cacheKey, user)
	if err != nil {
		dbUser, err := s.userRepo.GetUserByUsername(ctx, username)
		if err != nil {
			return nil, errors.Errorf("failed get user: %v", err)
		}
		user = dbUser
		_ = s.cache.SetStruct(ctx, cacheKey, user, cacheExpTime)
	}

	return user, nil
}

func (s *serv) getSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	session := &models.Session{}
	cacheKey := newUserCacheKey(sessionCacheKeyPrefix, refreshToken)
	err := s.cache.GetStruct(ctx, cacheKey, session)

	if err != nil {
		sessionDB, err := s.sessionRepo.GetSessionByToken(ctx, refreshToken)
		if err != nil {
			if errors.Is(err, repoerrs.ErrSessionNotExists) {
				return nil, serverrs.ErrSessionNotExists
			}

			return nil, fmt.Errorf("failed to get session: %w", err)
		}

		session = sessionDB
		_ = s.cache.SetStruct(ctx, cacheKey, session, cacheExpTime)
	}

	return session, nil
}

func (s *serv) createRefreshToken(userID int64, username, role string) (string, error) {
	refreshToken, err := utils.GenerateToken(
		models.UserInfo{
			UserID:   userID,
			Username: username,
			Role:     role,
		},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", errors.Errorf("failed to create refresh token: %v", err)
	}

	return refreshToken, nil
}

func (s *serv) createAccessToken(userID int64, username, role string) (string, error) {
	accessToken, err := utils.GenerateToken(
		models.UserInfo{
			UserID:   userID,
			Username: username,
			Role:     role,
		},
		[]byte(s.jwtConfig.AccessTokenSecret()),
		s.jwtConfig.AccessTokenExpiration(),
	)
	if err != nil {
		return "", errors.Errorf("failed to create access token: %v", err)
	}

	return accessToken, nil
}

func (s *serv) createSession(ctx context.Context, userID int64, refreshToken string) error {
	now := time.Now()
	session := &models.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		ExpiresAt:    now.Add(s.jwtConfig.RefreshTokenExpiration()),
		RevokedAt:    nil,
	}
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return errors.Errorf("failed to create a new session: %v", err)
	}

	return nil
}
