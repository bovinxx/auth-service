package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/pkg/errors"
)

func (s *Serv) GetAccessToken(
	ctx context.Context,
	token *models.Token,
) (*models.Token, error) {
	claims, err := utils.VerifyToken(token.Token, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return nil, errors.Errorf("failed to get refresh token: %v", err)
	}

	err = s.checkRefreshToken(ctx, token.Token)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.createAccessToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		Token: accessToken,
	}, nil
}
