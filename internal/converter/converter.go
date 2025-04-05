package converter

import (
	models "github.com/bovinxx/auth-service/internal/models/user"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
)

func ToCreateRequestFromRepo(user *models.User) *desc.CreateRequest {
	return &desc.CreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}

func ToGetResponseFromUser(user *models.User) *desc.GetResponse {
	return &desc.GetResponse{
		Id:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}
}

func ToUserFromCreateRequest(req *desc.CreateRequest) *models.User {
	return &models.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		IsAdmin:  req.GetIsAdmin(),
	}
}
