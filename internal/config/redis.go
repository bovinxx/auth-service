package config

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	redisHostEnvName              = "REDIS_HOST"
	redisPortEnvName              = "REDIS_PORT"
	redisPasswordEnvName          = "REDIS_PASSWORD"
	redisMaxIdleEnvName           = "REDIS_MAX_IDLE"
	redisIdleTimeoutEnvName       = "REDIS_IDLE_TIMEOUT"
	redisConnectionTimeoutEnvName = "REDIS_CONNECTION_TIMEOUT"
)

type RedisConfig interface {
	Address() string
	Password() string
	MaxIdle() int
	IdleTimeout() time.Duration
	ConnectionTimeout() time.Duration
}

type redis struct {
	address           string
	password          string
	maxIdle           int
	idleTimeout       time.Duration
	connectionTimeout time.Duration
}

func NewRedisConfig() (*redis, error) {
	host := os.Getenv(redisHostEnvName)
	if host == "" {
		return nil, errors.New("redis host not found")
	}

	port := os.Getenv(redisPortEnvName)
	if port == "" {
		return nil, errors.New("redis port not found")
	}

	password := os.Getenv(redisPasswordEnvName)
	if password == "" {
		return nil, errors.New("redis password not found")
	}

	maxIdleStr := os.Getenv(redisMaxIdleEnvName)
	if maxIdleStr == "" {
		return nil, errors.New("redis max idle not found")
	}

	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		return nil, errors.New("redis max idle is not an integer")
	}

	idleTimeoutStr := os.Getenv(redisIdleTimeoutEnvName)
	if idleTimeoutStr == "" {
		return nil, errors.New("redis idle timeout not found")
	}

	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		return nil, errors.New("redis idle timeout is not time duration")
	}

	connectionTimeoutStr := os.Getenv(redisConnectionTimeoutEnvName)
	if connectionTimeoutStr == "" {
		return nil, errors.New("redis connection timeout not found")
	}

	connectionTimeout, err := time.ParseDuration(connectionTimeoutStr)
	if err != nil {
		return nil, errors.New("redis connection timeout is not time duration")
	}

	return &redis{
		address:           net.JoinHostPort(host, port),
		password:          password,
		maxIdle:           maxIdle,
		idleTimeout:       idleTimeout,
		connectionTimeout: connectionTimeout,
	}, nil
}

func (cfg *redis) Address() string {
	return cfg.address
}

func (cfg *redis) Password() string {
	return cfg.password
}

func (cfg *redis) MaxIdle() int {
	return cfg.maxIdle
}

func (cfg *redis) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}

func (cfg *redis) ConnectionTimeout() time.Duration {
	return cfg.connectionTimeout
}
