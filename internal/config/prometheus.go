package config

import (
	"errors"
	"net"
	"os"
)

const (
	prometheusHostEnvName = "PROMETHEUS_HOST"
	prometheusPortEnvName = "PROMETHEUS_PORT"
)

type PrometheusConfig interface {
	Address() string
}

type prometheusConfig struct {
	address string
}

func NewPrometheusConfig() (*prometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if host == "" {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if port == "" {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		address: net.JoinHostPort(host, port),
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return cfg.address
}
