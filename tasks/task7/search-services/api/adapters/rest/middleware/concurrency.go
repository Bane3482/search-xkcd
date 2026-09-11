package middleware

import (
	"net/http"
)

type ConcurrencyLimiter struct {
	sem chan struct{}
}

func New(limit int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		sem: make(chan struct{}, limit),
	}
}

func (l *ConcurrencyLimiter) Acquire() bool {
	select {
	case l.sem <- struct{}{}:
		return true
	default:
		return false
	}

}

func (l *ConcurrencyLimiter) Release() {
	<-l.sem
}

func Concurrency(next http.HandlerFunc, limit int) http.HandlerFunc {
	limiter := New(limit)
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Acquire() {
			http.Error(w, "concurrency limit exceeded", http.StatusServiceUnavailable)
			return
		}
		defer limiter.Release()
		next(w, r)
	}
}
