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

func (s *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteUser(ctx, req.GetId())
	if err != nil {
		logger.Info("failed to delete user", logger.Err(err))
		if errors.Is(err, serviceerrs.ErrUserNotExists) {
			return nil, status.Errorf(codes.NotFound, "user not exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return nil, nil
}
