package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/pkg/errors"
)

func (s *serv) GetRefreshToken(
	ctx context.Context,
	token *models.Token,
) (*models.Token, error) {
	oldRefreshToken := token.Token
	claims, err := utils.VerifyToken(oldRefreshToken, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return nil, errors.Errorf("failed to get refresh token: %v", err)
	}

	username := claims.Username

	err = s.checkRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(
		models.UserInfo{
			UserID:   claims.UserID,
			Username: username,
			Role:     claims.Role,
		},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return nil, errors.Errorf("failed to create a new refresh token: %v", err)
	}

	session := &models.Session{
		UserID:       claims.UserID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(s.jwtConfig.RefreshTokenExpiration()),
		RevokedAt:    nil,
	}

	err = s.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new session: %w", err)
	}

	return &models.Token{
		Token: refreshToken,
	}, nil
}
