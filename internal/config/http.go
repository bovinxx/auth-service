package config

import (
	"errors"
	"fmt"
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
	host string
	port string
}

func NewHTTPConfig() (*httpConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return fmt.Sprintf("%s:%s", cfg.host, cfg.port)
}
