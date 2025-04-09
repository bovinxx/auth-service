package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
)

func (s *Implementation) GetAccessToken(
	ctx context.Context,
	req *desc.GetAccessTokenRequest,
) (*desc.GetAccessTokenResponse, error) {
	resp, err := s.service.GetAccessToken(ctx, converter.ToServiceFromGetAccessTokenRequest(req))
	if err != nil {
		return nil, err
	}
	return &desc.GetAccessTokenResponse{
		AccessToken: resp.Token,
	}, nil
}
