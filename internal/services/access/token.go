package access

import (
	"context"
	"strings"

	"github.com/bovinxx/auth-service/internal/models"
	"github.com/bovinxx/auth-service/internal/services/access/errors"
	"github.com/bovinxx/auth-service/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (s *serv) extractBearerToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.ErrNoMetadata
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.ErrNoAuthHeader
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", errors.ErrInvalidAuth
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	return accessToken, nil
}

func (s *serv) verifyToken(accessToken string) (*models.UserClaims, error) {
	claims, err := utils.VerifyToken(accessToken, []byte(s.jwtConfig.AccessTokenSecret()))
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
