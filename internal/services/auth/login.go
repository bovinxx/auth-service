package auth

import (
	"context"
	"time"

	"github.com/bovinxx/auth-service/internal/models"
	serverrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/pkg/errors"
)

func (s *serv) Login(ctx context.Context, req *models.User) (string, error) {
	user := &models.User{}
	cacheKey := newUserCacheKey(userCacheKeyPrefix, req.Name)

	err := s.cache.GetStruct(ctx, cacheKey, user)
	if err != nil {
		dbUser, err := s.userRepo.GetUserByUsername(ctx, req.Name)
		if err != nil {
			return "", errors.Errorf("failed get user: %v", err)
		}
		user = dbUser
		_ = s.cache.SetStruct(ctx, cacheKey, user, cacheExpTime)
	}

	if !utils.VerifyPassword(user.Password, req.Password) || (req.Role != user.Role) {
		return "", serverrs.ErrInvalidCredentials
	}

	refreshToken, err := utils.GenerateToken(
		models.UserInfo{
			UserID:   user.ID,
			Username: req.Name,
			Role:     req.Role,
		},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", errors.Errorf("failed to create refresh token: %v", err)
	}

	session := &models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(s.jwtConfig.RefreshTokenExpiration()),
		RevokedAt:    nil,
	}
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return "", errors.Errorf("failed to create a new session: %v", err)
	}

	return refreshToken, nil
}
