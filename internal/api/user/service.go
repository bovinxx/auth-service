package user

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

//go generate mockgen -source=service.go -destination=service_mock/mock.go -package=mock

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}

type Implementation struct {
	desc.UnimplementedUserServiceServer
	service UserService
}

func NewImplementation(service UserService) *Implementation {
	return &Implementation{
		service: service,
	}
}
