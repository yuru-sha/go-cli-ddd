// Package logger configures the application's default slog logger.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// InitLogger はアプリケーションのロガーを初期化します。
func InitLogger(logLevel string, debug bool) {
	level := getLogLevel(logLevel)

	var output io.Writer = os.Stdout
	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if debug {
		handler = slog.NewTextHandler(output, handlerOptions)
	} else {
		handler = slog.NewJSONHandler(output, handlerOptions)
	}

	slog.SetDefault(slog.New(handler))
}

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithAttrs は既定ロガーに属性を付与したロガーを返します。
func WithAttrs(args ...any) *slog.Logger {
	return slog.Default().With(args...)
}

// WithGroup は既定ロガーにグループを付与したロガーを返します。
func WithGroup(name string) *slog.Logger {
	return slog.Default().WithGroup(name)
}

// Default は既定ロガーを返します。
func Default() *slog.Logger {
	return slog.Default()
}

// FromContext は context に格納されたロガーを取得し、未設定なら既定ロガーを返します。
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok && logger != nil {
		return logger
	}
	return slog.Default()
}

// IntoContext はロガーを context に格納します。
func IntoContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

type loggerKey struct{}
