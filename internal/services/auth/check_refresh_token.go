package auth

import (
	"context"
	"time"

	serverrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
)

func (s *Serv) checkRefreshToken(ctx context.Context, refreshToken string) error {
	session, err := s.getSessionByToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	isSessionRevoked := session.RevokedAt != nil && !session.RevokedAt.IsZero()
	isSessionExpired := session.ExpiresAt.Before(time.Now())

	if isSessionRevoked {
		return serverrs.ErrSessionRevoked
	}

	if isSessionExpired {
		return serverrs.ErrSessionExpired
	}

	return nil
}
