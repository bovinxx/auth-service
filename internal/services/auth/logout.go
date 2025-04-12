package auth

import (
	"context"
	"fmt"

	"github.com/bovinxx/auth-service/internal/logger"
)

func (s *serv) Logout(ctx context.Context, refreshToken string) error {
	cacheKey := newUserCacheKey(sessionCacheKeyPrefix, refreshToken)
	err := s.cache.Del(ctx, cacheKey)
	if err != nil {
		logger.Warn("failed to delete session from cache", logger.Err(err))
	}
	err = s.sessionRepo.MarkSessionAsRevoked(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to mark session as revoked: %w", err)
	}
	return nil
}
