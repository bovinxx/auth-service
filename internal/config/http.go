package config

import (
	"errors"
	"net"
	"os"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	address string
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if host == "" {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if port == "" {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *httpConfig) Address() string {
	return cfg.address
}
