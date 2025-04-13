package converter

import (
	"github.com/bovinxx/auth-service/internal/models"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

func ToCreateRequestFromRepo(user *models.User) *desc.CreateRequest {
	return &desc.CreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     desc.Role(int32(user.Role)),
	}
}

func ToGetResponseFromUser(user *models.User) *desc.GetResponse {
	return &desc.GetResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  desc.Role(int32(user.Role)),
	}
}

func ToUserFromCreateRequest(req *desc.CreateRequest) *models.User {
	role := int32(req.GetRole())
	return &models.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     models.Role(role),
	}
}
