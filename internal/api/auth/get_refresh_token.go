package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	"github.com/bovinxx/auth-service/internal/logger"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	resp, err := s.service.GetRefreshToken(ctx, converter.ToServiceFromGetRefreshTokenRequest(req))
	if err != nil {
		logger.Info("failed to create a new refresh token", logger.Err(err))
		return nil, status.Errorf(codes.Internal, "failed to create a new refresh token: %v", err)
	}
	return &desc.GetRefreshTokenResponse{
		RefreshToken: resp.Token,
	}, nil
}
