package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/pkg/errors"
)

func (s *Serv) GetRefreshToken(
	ctx context.Context,
	token *models.Token,
) (*models.Token, error) {
	oldRefreshToken := token.Token
	claims, err := utils.VerifyToken(oldRefreshToken, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return nil, errors.Errorf("failed to get refresh token: %v", err)
	}

	err = s.checkRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return nil, err
	}

	err = s.createSession(ctx, claims.UserID, refreshToken)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		Token: refreshToken,
	}, nil
}
