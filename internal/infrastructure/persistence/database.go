package persistence

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// NewDatabase はデータベース接続を初期化します
func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// ロガーの設定
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// SQLiteデータベースに接続
	db, err := gorm.Open(sqlite.Open(cfg.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("データベース接続に失敗しました: %w", err)
	}

	log.Info().Str("dialect", cfg.Dialect).Str("dsn", cfg.DSN).Msg("データベースに接続しました")

	// データベースの初期化
	if cfg.AutoMigrate {
		log.Info().Msg("データベースマイグレーションを実行します")
		if err := InitDatabase(db); err != nil {
			return nil, fmt.Errorf("データベース初期化に失敗しました: %w", err)
		}
	}

	return db, nil
}
