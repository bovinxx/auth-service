package interceptor

import (
	"context"
	"time"

	ratelimiter "github.com/bovinxx/auth-service/pkg/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	rateLimit = 100
)

var (
	rateLimiter = ratelimiter.NewTokenBucketLimiter(context.Background(), rateLimit, time.Second)
)

func RateLimiterInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if rateLimiter.Allow() {
		return handler(ctx, req)
	}
	return nil, status.Error(codes.ResourceExhausted, "too many requests")
}
