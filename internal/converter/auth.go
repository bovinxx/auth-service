package converter

import (
	"github.com/bovinxx/auth-service/internal/models"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
)

func ToServiceFromLoginRequest(user *desc.LoginRequest) *models.User {
	return &models.User{
		Name:     user.GetUsername(),
		Password: user.GetPassword(),
		Role:     user.GetRole(),
	}
}

func ToServiceFromGetRefreshTokenRequest(req *desc.GetRefreshTokenRequest) *models.Token {
	return &models.Token{
		Token: req.GetOldRefreshToken(),
	}
}

func ToServiceFromGetAccessTokenRequest(req *desc.GetAccessTokenRequest) *models.Token {
	return &models.Token{
		Token: req.GetRefreshToken(),
	}
}

func ToLoginResponseFromService(req string) *desc.LoginResponse {
	return &desc.LoginResponse{
		RefreshToken: req,
	}
}
