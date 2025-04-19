package auth

import (
	"context"
	"errors"

	"github.com/bovinxx/auth-service/internal/converter"
	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
	desc "github.com/bovinxx/auth-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	resp, err := s.service.Login(ctx, converter.ToServiceFromLoginRequest(req))
	if err != nil {
		logger.Info("failed to login user", logger.Err(err))
		if errors.Is(err, serviceerrs.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Aborted, "invalid credentials")
		}
		if errors.Is(err, serviceerrs.ErrAccessDenied) {
			return nil, status.Errorf(codes.Aborted, "access denied")
		}
		return nil, status.Errorf(codes.Internal, "failed to create a new user: %v", err)
	}
	return converter.ToLoginResponseFromService(resp), nil
}
