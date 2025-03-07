package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// メインエントリーポイント
func main() {
	log.Info().Msg("GORM genによるモデル生成を開始します")

	// データベース接続
	db, err := connectDatabase()
	if err != nil {
		log.Error().Err(err).Msg("データベース接続に失敗しました")
		os.Exit(1)
	}

	// モデル生成
	if err := generateModels(db); err != nil {
		log.Error().Err(err).Msg("モデル生成に失敗しました")
		os.Exit(1)
	}

	log.Info().Msg("モデル生成が完了しました")
}

// connectDatabase はデータベースに接続します
func connectDatabase() (*gorm.DB, error) {
	// SQLiteデータベースに接続
	dsn := "file:go-cli-ddd.db?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // テーブル名を単数形にする
		},
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続に失敗しました: %w", err)
	}

	// テーブルが存在しない場合は作成
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables はテーブルを作成します
func createTables(db *gorm.DB) error {
	// テーブル定義
	type Campaign struct {
		ID        uint   `gorm:"primaryKey"`
		AccountID uint   `gorm:"not null;index"`
		Name      string `gorm:"size:255;not null"`
		Status    string `gorm:"size:50;not null"`
		Budget    int64  `gorm:"not null"`
		StartDate int64  `gorm:"column:start_date"`
		EndDate   int64  `gorm:"column:end_date"`
		CreatedAt int64  `gorm:"autoCreateTime"`
		UpdatedAt int64  `gorm:"autoUpdateTime"`
	}

	type Account struct {
		ID        uint   `gorm:"primaryKey"`
		Name      string `gorm:"size:255;not null"`
		Status    string `gorm:"size:50;not null"`
		APIKey    string `gorm:"column:api_key;size:255;not null"`
		CreatedAt int64  `gorm:"autoCreateTime"`
		UpdatedAt int64  `gorm:"autoUpdateTime"`

		Campaigns []Campaign `gorm:"foreignKey:AccountID;references:ID"`
	}

	// テーブル作成
	if err := db.AutoMigrate(&Account{}, &Campaign{}); err != nil {
		return fmt.Errorf("テーブル作成に失敗しました: %w", err)
	}

	return nil
}

// generateModels はGORM genを使用してモデルを生成します
func generateModels(db *gorm.DB) error {
	// ジェネレーターの設定
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/infrastructure/persistence/query",
		ModelPkgPath:      "./internal/infrastructure/persistence/model",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true, // NULL可能なフィールドにポインタを使用
		FieldCoverable:    true, // フィールドの上書きを許可
		FieldSignable:     true, // 符号付き整数を使用
		FieldWithIndexTag: true, // インデックスタグを生成
		FieldWithTypeTag:  true, // タイプタグを生成
	})

	// データベース接続の設定
	g.UseDB(db)

	// モデルの定義
	// Accountモデル
	account := g.GenerateModel("accounts",
		gen.FieldType("created_at", "time.Time"),
		gen.FieldType("updated_at", "time.Time"),
	)

	// Campaignモデル
	campaign := g.GenerateModel("campaigns",
		gen.FieldType("created_at", "time.Time"),
		gen.FieldType("updated_at", "time.Time"),
		gen.FieldType("start_date", "time.Time"),
		gen.FieldType("end_date", "time.Time"),
		// gen.FieldRelate(gen.RelateHasOne, "Account", account, &gen.FieldMeta{
		// 	FieldName: "AccountID",
		// }),
	)

	// 関連付け
	g.ApplyBasic(account, campaign)

	// コードの生成
	g.Execute()

	return nil
}
