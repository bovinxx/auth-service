package access

import (
	"context"
	"strings"

	"github.com/bovinxx/auth-service/internal/services/access/errors"
	"google.golang.org/grpc/metadata"
)

func extractBearerToken(ctx context.Context) (string, error) {
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
