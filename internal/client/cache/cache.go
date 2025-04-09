package cache

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	GetStruct(ctx context.Context, key string, out interface{}) error
	SetStruct(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Ping(ctx context.Context) error
}
