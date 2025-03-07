package persistence

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// GenerateModels はGORM genを使用してモデルを生成します
func GenerateModels(db *gorm.DB) error {
	log.Info().Msg("モデルの生成を開始します")

	// ジェネレーターの設定
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/infrastructure/persistence/query",
		ModelPkgPath: "./internal/infrastructure/persistence/model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// データベース接続の設定
	g.UseDB(db)

	// モデルの定義
	g.ApplyBasic(
		// Account モデル
		g.GenerateModel("accounts", gen.FieldType("created_at", "time.Time"), gen.FieldType("updated_at", "time.Time")),
		
		// Campaign モデル
		g.GenerateModel("campaigns", 
			gen.FieldType("created_at", "time.Time"), 
			gen.FieldType("updated_at", "time.Time"),
			gen.FieldType("start_date", "time.Time"),
			gen.FieldType("end_date", "time.Time"),
		),
	)

	// コードの生成
	g.Execute()

	log.Info().Msg("モデルの生成が完了しました")
	return nil
}

// InitDatabase はデータベースの初期化とマイグレーションを行います
func InitDatabase(db *gorm.DB) error {
	// テーブルの作成
	if err := db.AutoMigrate(&entity.Account{}, &entity.Campaign{}); err != nil {
		return err
	}

	// モデルの生成
	if err := GenerateModels(db); err != nil {
		return err
	}

	return nil
}
