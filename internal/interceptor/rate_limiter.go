package interceptor

import (
	"context"
	"time"

	ratelimiter "github.com/bovinxx/auth-service/internal/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	RATE_LIMIT = 100
)

var (
	rateLimiter = ratelimiter.NewTokenBucketLimiter(context.Background(), RATE_LIMIT, time.Second)
)

func RateLimiterInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if rateLimiter.Allow() {
		return handler(ctx, req)
	} else {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}
}
