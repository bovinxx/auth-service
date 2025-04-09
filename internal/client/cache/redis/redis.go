package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bovinxx/auth-service/internal/client/cache"
	"github.com/bovinxx/auth-service/internal/config"
	"github.com/gomodule/redigo/redis"
)

var _ cache.RedisClient = (*client)(nil)

type client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

func NewClient(pool *redis.Pool, config config.RedisConfig) cache.RedisClient {
	return &client{
		pool:   pool,
		config: config,
	}
}

func (c *client) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, int(ttl.Seconds()), value)
	return err
}

func (c *client) Get(_ context.Context, key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", fmt.Errorf("key not found")
	}

	return reply, err
}

func (c *client) Del(_ context.Context, key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (c *client) GetStruct(ctx context.Context, key string, out interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), out)
	if err != nil {
		return fmt.Errorf("failed to get struct: %w", err)
	}
	return nil
}

func (c *client) SetStruct(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to set struct: %w", err)
	}

	return c.Set(ctx, key, string(data), ttl)
}

func (c *client) Ping(_ context.Context) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err
}

func (c *client) Close() error {
	return c.pool.Close()
}
