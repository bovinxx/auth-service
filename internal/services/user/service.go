package service

import (
	"context"

	models "github.com/bovinxx/auth-service/internal/models/user"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, name, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}
