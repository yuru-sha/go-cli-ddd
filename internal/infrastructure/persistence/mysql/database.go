package mysql

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// Database はデータベース接続を管理します
type Database struct {
	DB     *gorm.DB
	Config *config.Config
}

// RandomPolicy はランダムなレプリカを選択するポリシーです
type RandomPolicy struct {
}

// Resolve はランダムなレプリカを選択します
func (p RandomPolicy) Resolve(replicas []gorm.ConnPool) gorm.ConnPool {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(replicas))))
	if err != nil {
		// エラーが発生した場合は最初のレプリカを返す
		log.Error().Err(err).Msg("ランダムな数値の生成に失敗しました")
		return replicas[0]
	}
	return replicas[n.Int64()]
}

// RoundRobinPolicy はラウンドロビン方式でレプリカを選択するポリシーです
type RoundRobinPolicy struct {
	counter int
	mu      sync.Mutex
}

// Resolve はラウンドロビン方式でレプリカを選択します
func (p *RoundRobinPolicy) Resolve(replicas []gorm.ConnPool) gorm.ConnPool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter = (p.counter + 1) % len(replicas)
	return replicas[p.counter]
}

// NewDatabase は新しいデータベース接続を作成します
func NewDatabase(cfg *config.Config) (*Database, error) {
	// ログレベルを設定
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

	// GORMの設定
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// データベースインスタンス
	db := &Database{
		Config: cfg,
	}

	// Aurora接続が有効な場合
	if cfg.Database.Aurora.Enabled {
		// Auroraクラスター接続を設定
		gormDB, err := db.setupAuroraConnection(context.Background(), gormConfig)
		if err != nil {
			return nil, err
		}
		db.DB = gormDB
	} else {
		// 通常の単一接続を設定
		dsn, err := getDSN(context.Background(), cfg)
		if err != nil {
			return nil, fmt.Errorf("DSNの取得に失敗しました: %w", err)
		}

		// データベースに接続
		var gormDB *gorm.DB
		var dbErr error
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

	// 自動マイグレーションが有効な場合は実行
	if cfg.Database.AutoMigrate {
		if err := autoMigrate(db.DB); err != nil {
			return nil, fmt.Errorf("マイグレーションに失敗しました: %w", err)
		}
	}

	return db, nil
}

// setupAuroraConnection はAuroraクラスター接続を設定します
func (db *Database) setupAuroraConnection(ctx context.Context, gormConfig *gorm.Config) (*gorm.DB, error) {
	cfg := db.Config

	// Secret Managerクライアントを作成
	secretsManager, err := secrets.NewAWSSecretsManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("Secret Managerの初期化に失敗しました: %w", err)
	}

	// ライター接続情報を取得
	if cfg.Database.Aurora.Writer.SecretID == "" {
		return nil, fmt.Errorf("Auroraライター接続のSecretIDが設定されていません")
	}

	writerSecret, err := secretsManager.GetDatabaseSecret(ctx, cfg.Database.Aurora.Writer.SecretID)
	if err != nil {
		return nil, fmt.Errorf("Auroraライター接続情報の取得に失敗しました: %w", err)
	}

	// ライター接続を作成
	writerDSN := writerSecret.FormatDSN(cfg.Database.Dialect)
	var gormDB *gorm.DB

	switch cfg.Database.Dialect {
	case "mysql":
		gormDB, err = gorm.Open(mysql.Open(writerDSN), gormConfig)
	case "postgres":
		gormDB, err = gorm.Open(postgres.Open(writerDSN), gormConfig)
	default:
		return nil, fmt.Errorf("Auroraでは未対応のデータベースダイアレクト: %s", cfg.Database.Dialect)
	}

	if err != nil {
		return nil, fmt.Errorf("Auroraライター接続に失敗しました: %w", err)
	}

	// リーダー接続情報を取得
	if cfg.Database.Aurora.Reader.SecretID != "" {
		readerSecret, err := secretsManager.GetDatabaseSecret(ctx, cfg.Database.Aurora.Reader.SecretID)
		if err != nil {
			return nil, fmt.Errorf("Auroraリーダー接続情報の取得に失敗しました: %w", err)
		}

		// リーダー接続を設定
		readerDSN := readerSecret.FormatDSN(cfg.Database.Dialect)

		// DBResolverを使用してリーダー/ライターを設定
		resolverConfig := dbresolver.Config{
			Replicas: []gorm.Dialector{},
			Policy:   db.getLoadBalancingPolicy(),
		}

		// リーダーダイアレクタを作成
		var readerDialector gorm.Dialector
		switch cfg.Database.Dialect {
		case "mysql":
			readerDialector = mysql.Open(readerDSN)
		case "postgres":
			readerDialector = postgres.Open(readerDSN)
		}

		resolverConfig.Replicas = append(resolverConfig.Replicas, readerDialector)

		// DBResolverを登録
		err = gormDB.Use(dbresolver.Register(resolverConfig).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(10).
			SetMaxOpenConns(100))

		if err != nil {
			return nil, fmt.Errorf("DBResolverの設定に失敗しました: %w", err)
		}

		log.Info().Msg("Auroraクラスター接続（ライター/リーダー）を設定しました")
	} else {
		log.Info().Msg("Auroraライター接続のみを設定しました（リーダーは設定されていません）")
	}

	return gormDB, nil
}

// getLoadBalancingPolicy はロードバランシングポリシーを取得します
func (db *Database) getLoadBalancingPolicy() dbresolver.Policy {
	switch db.Config.Database.Aurora.Reader.LoadBalancing {
	case "round-robin":
		return &RoundRobinPolicy{}
	case "random":
		return &RandomPolicy{}
	default:
		return &RandomPolicy{} // デフォルトはランダム
	}
}

// getDSN はデータベース接続文字列を取得します
// Secret Managerが有効な場合はそこから取得し、そうでない場合は設定から取得します
func getDSN(ctx context.Context, cfg *config.Config) (string, error) {
	// Secret Managerが有効で、SecretIDが設定されている場合
	if cfg.AWS.Secrets.Enabled && cfg.Database.SecretID != "" {
		log.Info().Msg("Secret Managerからデータベース接続情報を取得します")

		// Secret Managerクライアントを作成
		secretsManager, err := secrets.NewAWSSecretsManager(cfg)
		if err != nil {
			return "", fmt.Errorf("Secret Managerの初期化に失敗しました: %w", err)
		}

		// データベース接続情報を取得
		dbSecret, err := secretsManager.GetDatabaseSecret(ctx, cfg.Database.SecretID)
		if err != nil {
			return "", fmt.Errorf("データベース接続情報の取得に失敗しました: %w", err)
		}

		// 接続文字列を生成
		return dbSecret.FormatDSN(cfg.Database.Dialect), nil
	}

	// Secret Managerが無効または設定されていない場合は設定から取得
	log.Info().Msg("設定ファイルからデータベース接続情報を使用します")
	return cfg.Database.DSN, nil
}

// GetReader はリードオンリー操作用のデータベース接続を返します
func (db *Database) GetReader() *gorm.DB {
	// DBResolverを使用している場合、Clauses(dbresolver.Read)を使用してリードレプリカに接続
	return db.DB.Clauses(dbresolver.Read)
}

// GetWriter はライト操作用のデータベース接続を返します
func (db *Database) GetWriter() *gorm.DB {
	// DBResolverを使用している場合、Clauses(dbresolver.Write)を使用してライターに接続
	return db.DB.Clauses(dbresolver.Write)
}

// autoMigrate はデータベースのマイグレーションを実行します
func autoMigrate(_ *gorm.DB) error {
	// ここにマイグレーション対象のモデルを追加
	// 例: return db.AutoMigrate(&models.User{}, &models.Task{})
	return nil
}
