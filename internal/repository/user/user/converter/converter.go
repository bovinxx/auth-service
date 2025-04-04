package converter

import (
	models "github.com/bovinxx/auth-service/internal/models/user"
	"github.com/bovinxx/auth-service/internal/repository/user/user/model"
)

func ToUserFromRepo(user *model.User) *models.User {
	return &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}
