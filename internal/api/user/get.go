package user

import (
	"context"
	"errors"

	"github.com/bovinxx/auth-service/internal/converter"
	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/user/errors"
	desc "github.com/bovinxx/auth-service/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.service.GetUser(ctx, req.GetId())
	if err != nil {
		logger.Info("failed to get user", logger.Err(err))
		if errors.Is(err, serviceerrs.ErrUserNotExists) {
			return nil, status.Errorf(codes.NotFound, "user not exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return converter.ToGetResponseFromUser(user), nil
}
