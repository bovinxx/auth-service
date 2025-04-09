package config

import (
	"errors"
	"os"
	"time"
)

const (
	jwtAccessTokenSecretEnvName      = "ACCESS_TOKEN_SECRET"
	jwtAccessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"
	jwtRefreshTokenSecretEnvName     = "REFRESH_TOKEN_SECRET"
	jwtRefreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION"
)

type JWTConfig interface {
	AccessTokenSecret() string
	AccessTokenExpiration() time.Duration
	RefreshTokenSecret() string
	RefreshTokenExpiration() time.Duration
}

type jwtConfig struct {
	accessTokenSecret      string
	accessTokenExpiration  time.Duration
	refreshTokenSecret     string
	refreshTokenExpiration time.Duration
}

func NewJWTConfig() (JWTConfig, error) {
	accessSecret := os.Getenv(jwtAccessTokenSecretEnvName)
	if accessSecret == "" {
		return nil, errors.New("access secret not found")
	}

	refreshSecret := os.Getenv(jwtRefreshTokenSecretEnvName)
	if refreshSecret == "" {
		return nil, errors.New("refresh secret not found")
	}

	accessExpStr := os.Getenv(jwtAccessTokenExpirationEnvName)
	accessExp, err := time.ParseDuration(accessExpStr)
	if err != nil {
		return nil, errors.New("failed to parse access expiration")
	}

	refreshExpStr := os.Getenv(jwtRefreshTokenExpirationEnvName)
	refreshExp, err := time.ParseDuration(refreshExpStr)
	if err != nil {
		return nil, errors.New("failed to parse refresh expiration")
	}

	return &jwtConfig{
		accessTokenSecret:      accessSecret,
		accessTokenExpiration:  accessExp,
		refreshTokenSecret:     refreshSecret,
		refreshTokenExpiration: refreshExp,
	}, nil
}

func (cfg *jwtConfig) AccessTokenSecret() string {
	return cfg.accessTokenSecret
}

func (cfg *jwtConfig) AccessTokenExpiration() time.Duration {
	return cfg.accessTokenExpiration
}

func (cfg *jwtConfig) RefreshTokenSecret() string {
	return cfg.refreshTokenSecret
}

func (cfg *jwtConfig) RefreshTokenExpiration() time.Duration {
	return cfg.refreshTokenExpiration
}
