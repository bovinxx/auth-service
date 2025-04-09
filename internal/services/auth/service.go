package service

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
)

type Service interface {
	Login(ctx context.Context, req *models.User) (string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetRefreshToken(ctx context.Context, token *models.Token) (*models.Token, error)
	GetAccessToken(ctx context.Context, token *models.Token) (*models.Token, error)
}
