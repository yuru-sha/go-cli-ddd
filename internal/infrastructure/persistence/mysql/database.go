package mysql

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// Database はデータベース接続を管理します。
type Database struct {
	DB     *gorm.DB
	Config *config.Config
}

// RandomPolicy picks one replica at random.
type RandomPolicy struct{}

// Resolve chooses one replica connection pool.
func (p RandomPolicy) Resolve(replicas []gorm.ConnPool) gorm.ConnPool {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(replicas))))
	if err != nil {
		slog.Error("ランダムな数値の生成に失敗しました", "err", err)
		return replicas[0]
	}
	return replicas[n.Int64()]
}

// RoundRobinPolicy rotates through replicas in order.
type RoundRobinPolicy struct {
	counter int
	mu      sync.Mutex
}

// Resolve chooses one replica connection pool using round-robin order.
func (p *RoundRobinPolicy) Resolve(replicas []gorm.ConnPool) gorm.ConnPool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter = (p.counter + 1) % len(replicas)
	return replicas[p.counter]
}

// NewDatabase creates a configured database wrapper.
func NewDatabase(cfg *config.Config) (*Database, error) {
	var logLevel logger.LogLevel
	switch cfg.Database.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}

	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	db := &Database{Config: cfg}

	if cfg.Database.Aurora.Enabled {
		gormDB, err := db.setupAuroraConnection(context.Background(), gormConfig)
		if err != nil {
			return nil, err
		}
		db.DB = gormDB
	} else {
		dsn, err := getDSN(context.Background(), cfg)
		if err != nil {
			return nil, fmt.Errorf("DSNの取得に失敗しました: %w", err)
		}

		var (
			gormDB *gorm.DB
			dbErr  error
		)
		switch cfg.Database.Dialect {
		case "sqlite":
			gormDB, dbErr = gorm.Open(sqlite.Open(dsn), gormConfig)
		case "mysql":
			gormDB, dbErr = gorm.Open(mysql.Open(dsn), gormConfig)
		case "postgres":
			gormDB, dbErr = gorm.Open(postgres.Open(dsn), gormConfig)
		default:
			return nil, fmt.Errorf("未対応のデータベースダイアレクト: %s", cfg.Database.Dialect)
		}

		if dbErr != nil {
			return nil, fmt.Errorf("データベース接続に失敗しました: %w", dbErr)
		}

		db.DB = gormDB
	}

	if cfg.Database.AutoMigrate {
		if err := autoMigrate(db.DB); err != nil {
			return nil, fmt.Errorf("マイグレーションに失敗しました: %w", err)
		}
	}

	return db, nil
}

func (db *Database) setupAuroraConnection(ctx context.Context, gormConfig *gorm.Config) (*gorm.DB, error) {
	cfg := db.Config

	if !cfg.UseAWSSecretsManager() {
		return nil, fmt.Errorf("aurora接続は secrets provider=aws のときだけ利用できます")
	}

	secretsManager, err := secrets.NewAWSSecretsManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("secret managerの初期化に失敗しました: %w", err)
	}

	if cfg.Database.Aurora.Writer.SecretID == "" {
		return nil, fmt.Errorf("Auroraライター接続のSecretIDが設定されていません")
	}

	writerSecret, err := secretsManager.GetDatabaseSecret(ctx, cfg.Database.Aurora.Writer.SecretID)
	if err != nil {
		return nil, fmt.Errorf("auroraライター接続情報の取得に失敗しました: %w", err)
	}

	writerDSN := writerSecret.FormatDSN(cfg.Database.Dialect)
	var gormDB *gorm.DB

	switch cfg.Database.Dialect {
	case "mysql":
		gormDB, err = gorm.Open(mysql.Open(writerDSN), gormConfig)
	case "postgres":
		gormDB, err = gorm.Open(postgres.Open(writerDSN), gormConfig)
	default:
		return nil, fmt.Errorf("auroraでは未対応のデータベースダイアレクト: %s", cfg.Database.Dialect)
	}

	if err != nil {
		return nil, fmt.Errorf("auroraライター接続に失敗しました: %w", err)
	}

	if cfg.Database.Aurora.Reader.SecretID != "" {
		readerSecret, err := secretsManager.GetDatabaseSecret(ctx, cfg.Database.Aurora.Reader.SecretID)
		if err != nil {
			return nil, fmt.Errorf("auroraリーダー接続情報の取得に失敗しました: %w", err)
		}

		readerDSN := readerSecret.FormatDSN(cfg.Database.Dialect)
		resolverConfig := dbresolver.Config{
			Replicas: []gorm.Dialector{},
			Policy:   db.getLoadBalancingPolicy(),
		}

		var readerDialector gorm.Dialector
		switch cfg.Database.Dialect {
		case "mysql":
			readerDialector = mysql.Open(readerDSN)
		case "postgres":
			readerDialector = postgres.Open(readerDSN)
		}

		resolverConfig.Replicas = append(resolverConfig.Replicas, readerDialector)

		err = gormDB.Use(dbresolver.Register(resolverConfig).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(10).
			SetMaxOpenConns(100))
		if err != nil {
			return nil, fmt.Errorf("DBResolverの設定に失敗しました: %w", err)
		}

		slog.Info("Auroraクラスター接続（ライター/リーダー）を設定しました")
	} else {
		slog.Info("Auroraライター接続のみを設定しました（リーダーは設定されていません）")
	}

	return gormDB, nil
}

func (db *Database) getLoadBalancingPolicy() dbresolver.Policy {
	switch db.Config.Database.Aurora.Reader.LoadBalancing {
	case "round-robin":
		return &RoundRobinPolicy{}
	case "random":
		return &RandomPolicy{}
	default:
		return &RandomPolicy{}
	}
}

func getDSN(ctx context.Context, cfg *config.Config) (string, error) {
	if cfg.UseAWSSecretsManager() && cfg.Database.SecretID != "" {
		slog.Info("Secret Managerからデータベース接続情報を取得します")

		secretsManager, err := secrets.NewAWSSecretsManager(cfg)
		if err != nil {
			return "", fmt.Errorf("secret managerの初期化に失敗しました: %w", err)
		}

		dbSecret, err := secretsManager.GetDatabaseSecret(ctx, cfg.Database.SecretID)
		if err != nil {
			return "", fmt.Errorf("データベース接続情報の取得に失敗しました: %w", err)
		}

		return dbSecret.FormatDSN(cfg.Database.Dialect), nil
	}

	slog.Info("設定ファイルからデータベース接続情報を使用します")
	return cfg.Database.DSN, nil
}

// GetReader returns the read-side database handle.
func (db *Database) GetReader() *gorm.DB {
	return db.DB.Clauses(dbresolver.Read)
}

// GetWriter returns the write-side database handle.
func (db *Database) GetWriter() *gorm.DB {
	return db.DB.Clauses(dbresolver.Write)
}

func autoMigrate(_ *gorm.DB) error {
	return nil
}
