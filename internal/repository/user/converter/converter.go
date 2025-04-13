package converter

import (
	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/repository/user/model"
)

func ToUserFromRepo(user *model.User) *models.User {
	return &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     models.Role(user.Role),
	}
}
