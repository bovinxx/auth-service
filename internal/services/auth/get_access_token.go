package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/utils"
	"github.com/pkg/errors"
)

func (s *serv) GetAccessToken(
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

	accessToken, err := utils.GenerateToken(
		models.UserInfo{
			Username: claims.Username,
			Role:     claims.Role,
		},
		[]byte(s.jwtConfig.AccessTokenSecret()),
		s.jwtConfig.AccessTokenExpiration(),
	)
	if err != nil {
		return nil, errors.Errorf("failed to create a new access token: %v", err)
	}

	return &models.Token{
		Token: accessToken,
	}, nil
}
