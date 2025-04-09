package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func AuthzInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	return nil, nil
}
