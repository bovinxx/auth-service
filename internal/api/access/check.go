package access

import (
	"context"
	"errors"

	"github.com/bovinxx/auth-service/internal/logger"
	serviceerrs "github.com/bovinxx/auth-service/internal/services/access/errors"
	desc "github.com/bovinxx/auth-service/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	_, err := s.service.Check(ctx, req.GetEndpointAddress())
	if err == nil {
		return &emptypb.Empty{}, nil
	}

	logger.Info("failed to check access", logger.Err(err))

	switch {
	case errors.Is(err, serviceerrs.ErrAccessDenied):
		return nil, status.Error(codes.PermissionDenied, "access denied")
	case errors.Is(err, serviceerrs.ErrNoAuthHeader):
		return nil, status.Error(codes.Unauthenticated, "no auth header")
	case errors.Is(err, serviceerrs.ErrInvalidToken):
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	case errors.Is(err, serviceerrs.ErrInvalidAuth):
		return nil, status.Error(codes.PermissionDenied, "invalid credentials")
	case errors.Is(err, serviceerrs.ErrNoMetadata):
		return nil, status.Error(codes.FailedPrecondition, "no metadata found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to check access: %v", err)
	}
}
