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

	// 整数オーバーフローを防ぐため、maxRetriesが負の値や大きすぎる値の場合は制限する
	if maxRetries < 0 {
		maxRetries = 0
	}

	// uint64の最大値を超えないようにする
	// 安全な値として100を上限とする
	const maxSafeRetries = 100
	if maxRetries > maxSafeRetries {
		maxRetries = maxSafeRetries
	}

	// 安全に変換
	var maxRetriesUint64 uint64
	if maxRetries >= 0 {
		maxRetriesUint64 = uint64(maxRetries)
	} else {
		maxRetriesUint64 = 0
	}

	return backoff.WithMaxRetries(exponentialBackOff, maxRetriesUint64)
}
