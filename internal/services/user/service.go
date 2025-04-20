package user

import (
	"context"

	"github.com/bovinxx/auth-service/internal/models"
)

// go generate mockgen -source=service.go -destination=service_mock/mock.go -package=mock

type userRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}

type Serv struct {
	repo userRepository
}

func NewService(repo userRepository) *Serv {
	return &Serv{
		repo: repo,
	}
}
