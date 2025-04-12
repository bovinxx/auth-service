package user

import (
	"context"

	"github.com/bovinxx/auth-service/internal/client/db"
	"github.com/bovinxx/auth-service/internal/models"
)

type userRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}

type serv struct {
	repo      userRepository
	txManager db.TxManager
}

func NewService(repo userRepository, txManager db.TxManager) *serv {
	return &serv{
		repo:      repo,
		txManager: txManager,
	}
}
