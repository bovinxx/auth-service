package auth

import (
	"context"

	"github.com/bovinxx/auth-service/internal/converter"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
)

func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	resp, err := s.service.Login(ctx, converter.ToServiceFromLoginRequest(req))
	if err != nil {
		return nil, err
	}
	return converter.ToLoginResponseFromService(resp), nil
}
