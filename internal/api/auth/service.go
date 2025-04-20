package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
)

// go generate mockgen -source=service.go -destination=service_mock/mock.go -package=mock

type Service interface {
	Login(ctx context.Context, req *models.User) (string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetRefreshToken(ctx context.Context, token *models.Token) (*models.Token, error)
	GetAccessToken(ctx context.Context, token *models.Token) (*models.Token, error)
}

type Implementation struct {
	desc.UnimplementedAuthServiceServer
	service Service
}

func NewImplementation(service Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
