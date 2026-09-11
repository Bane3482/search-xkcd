package middleware

import (
	"context"
	"net/http"
	"time"
)

type TokenBucketLimiter struct {
	tokenBucketCh chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, limit int, period time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokenBucketCh: make(chan struct{}, 1),
	}

	interval := period.Nanoseconds() / int64(limit)

	go limiter.startPeriodReplenishment(ctx, time.Duration(interval))

	return limiter
}

func (l *TokenBucketLimiter) startPeriodReplenishment(ctx context.Context, interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			select {
			case l.tokenBucketCh <- struct{}{}:
			default:
			}
		}
	}
}

func (l *TokenBucketLimiter) Allow(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-l.tokenBucketCh:
		return
	}
}

func Rate(next http.HandlerFunc, rps int) http.HandlerFunc {
	limiter := NewTokenBucketLimiter(context.Background(), rps, time.Second)
	return func(w http.ResponseWriter, r *http.Request) {
		limiter.Allow(r.Context())
		next(w, r)
	}
}
