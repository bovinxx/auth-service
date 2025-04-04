package user

import (
	"context"

	models "github.com/bovinxx/auth-service/internal/models/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, name, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id int64) error
}
