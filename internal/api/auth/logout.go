package auth

import (
	"context"

	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Logout(ctx context.Context, req *desc.LogoutRequest) (*emptypb.Empty, error) {
	err := s.service.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
