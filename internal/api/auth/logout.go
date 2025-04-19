package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/logger"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Logout(ctx context.Context, req *desc.LogoutRequest) (*emptypb.Empty, error) {
	err := s.service.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		logger.Info("failed to logout", logger.Err(err))
		return nil, status.Errorf(codes.Internal, "failed to logout: %v", err)
	}
	return &emptypb.Empty{}, nil
}
