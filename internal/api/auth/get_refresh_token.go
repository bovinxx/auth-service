package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
)

func (s *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	resp, err := s.service.GetRefreshToken(ctx, converter.ToServiceFromGetRefreshTokenRequest(req))
	if err != nil {
		return nil, err
	}
	return &desc.GetRefreshTokenResponse{
		RefreshToken: resp.Token,
	}, nil
}
