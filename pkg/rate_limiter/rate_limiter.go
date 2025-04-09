package ratelimiter

import (
	"context"
	"time"
)

type TokenBucketLimiter struct {
	tokenBucketCh chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, limit int, period time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokenBucketCh: make(chan struct{}, limit),
	}

	for i := 0; i < limit; i++ {
		limiter.tokenBucketCh <- struct{}{}
	}

	replenishmentInterval := period.Nanoseconds() / int64(limit)
	go limiter.startPeriodReplenishment(ctx, replenishmentInterval)

	return limiter
}

func (limiter *TokenBucketLimiter) startPeriodReplenishment(ctx context.Context, replenishmentInterval int64) {
	ticker := time.NewTicker(time.Duration(replenishmentInterval))

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			limiter.tokenBucketCh <- struct{}{}
		}
	}
}

func (limiter *TokenBucketLimiter) Allow() bool {
	select {
	case <-limiter.tokenBucketCh:
		return true
	default:
		return false
	}
}
