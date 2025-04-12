package interceptor

import (
	"context"

	"github.com/bovinxx/auth-service/internal/api/access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthzInterceptor(accessSvc access.AccessService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		ok, err := accessSvc.Check(ctx, info.FullMethod)
		if err != nil || !ok {
			return nil, status.Error(codes.PermissionDenied, "not enough rights")
		}

		return handler(ctx, req)
	}
}
