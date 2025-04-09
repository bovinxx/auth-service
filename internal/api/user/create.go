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

func (s *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.service.CreateUser(ctx, converter.ToUserFromCreateRequest(req))
	if err != nil {
		logger.Info("failed to create a new user", logger.Err(err))
		if errors.Is(err, serviceerrs.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create a new user: %v", err)
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
