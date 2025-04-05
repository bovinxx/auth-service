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

// type AuthService interface {
// 	Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error)
// 	GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error)
// 	GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error)
// }

// type AccessService interface {
// 	Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error)
// }
