package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger はアプリケーションのロガーを初期化します
func InitLogger(logLevel string, debug bool) {
	// ログレベルの設定
	level := getLogLevel(logLevel)
	zerolog.SetGlobalLevel(level)

	// 開発モードの場合はより読みやすい出力形式を使用
	var output io.Writer = os.Stdout
	if debug {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// グローバルロガーの設定
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

// getLogLevel は文字列のログレベルをzerolog.Levelに変換します
func getLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
