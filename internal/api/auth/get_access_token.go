package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	"github.com/bovinxx/auth-service/internal/logger"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetAccessToken(
	ctx context.Context,
	req *desc.GetAccessTokenRequest,
) (*desc.GetAccessTokenResponse, error) {
	resp, err := s.service.GetAccessToken(ctx, converter.ToServiceFromGetAccessTokenRequest(req))
	if err != nil {
		logger.Info("failed to create a new access token", logger.Err(err))
		return nil, status.Errorf(codes.Internal, "failed to create a new access token: %v", err)
	}
	return &desc.GetAccessTokenResponse{
		AccessToken: resp.Token,
	}, nil
}
