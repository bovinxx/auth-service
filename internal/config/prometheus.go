package config

import (
	"errors"
	"fmt"
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
	host string
	port string
}

func NewPrometheusConfig() (*prometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return fmt.Sprintf("%s:%s", cfg.host, cfg.port)
}
