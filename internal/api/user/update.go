package user

import (
	"context"
	"errors"

	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateUser(
		ctx,
		req.GetId(),
		req.GetOldPassword(),
		req.GetNewPassword(),
	)
	if err != nil {
		logger.Info("failed to update user", logger.Err(err))
		if errors.Is(err, serviceerrs.ErrUserNotExists) {
			return nil, status.Errorf(codes.NotFound, "user not exist")
		}
		if errors.Is(err, serviceerrs.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return nil, nil
}
