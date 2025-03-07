package http

import (
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/time/rate"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// NewHTTPClient は設定に基づいてHTTPクライアントを作成します
func NewHTTPClient(cfg *config.HTTPConfig) *http.Client {
	return &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// NewRateLimiter はレート制限を行うリミッターを作成します
func NewRateLimiter(cfg *config.RateLimitConfig) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(cfg.QPS), cfg.Burst)
}

// NewBackOff はリトライ用のバックオフポリシーを作成します
func NewBackOff(maxRetries int) backoff.BackOff {
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.InitialInterval = 100 * time.Millisecond
	exponentialBackOff.MaxInterval = 10 * time.Second
	exponentialBackOff.MaxElapsedTime = 30 * time.Second

	return backoff.WithMaxRetries(exponentialBackOff, uint64(maxRetries))
}
