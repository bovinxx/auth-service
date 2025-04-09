package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/bovinxx/auth-service/internal/models"
	repoerrs "github.com/bovinxx/auth-service/internal/repository/session/errors"
	serverrs "github.com/bovinxx/auth-service/internal/services/auth/errors"
	"github.com/pkg/errors"
)

func (s *serv) checkRefreshToken(ctx context.Context, refreshToken string) error {
	session := &models.Session{}
	cacheKey := newUserCacheKey(sessionCacheKeyPrefix, refreshToken)
	err := s.cache.GetStruct(ctx, cacheKey, session)

	if err != nil {
		sessionDB, err := s.sessionRepo.GetSessionByToken(ctx, refreshToken)
		if err != nil {
			if errors.Is(err, repoerrs.ErrSessionNotExists) {
				return serverrs.ErrSessionNotExists
			}

			return fmt.Errorf("failed to get session: %w", err)
		}

		session = sessionDB
		_ = s.cache.SetStruct(ctx, cacheKey, session, cacheExpTime)
	}
	if session.RevokedAt != nil && !session.RevokedAt.IsZero() {
		return serverrs.ErrSessionRevoked
	}

	if session.ExpiresAt.Before(time.Now()) {
		return serverrs.ErrSessionExpired
	}

	return nil
}
