package mysql

import (
	"log/slog"

	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// GenerateModels はGORM genを使用してモデルを生成します。
func GenerateModels(db *gorm.DB) error {
	slog.Info("モデルの生成を開始します")

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/infrastructure/persistence/query",
		ModelPkgPath: "./internal/infrastructure/persistence/model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("accounts", gen.FieldType("created_at", "time.Time"), gen.FieldType("updated_at", "time.Time")),
		g.GenerateModel("campaigns",
			gen.FieldType("created_at", "time.Time"),
			gen.FieldType("updated_at", "time.Time"),
			gen.FieldType("start_date", "time.Time"),
			gen.FieldType("end_date", "time.Time"),
		),
	)

	g.Execute()

	slog.Info("モデルの生成が完了しました")
	return nil
}

// InitDatabase はデータベースの初期化とマイグレーションを行います。
func InitDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Account{}, &entity.Campaign{}); err != nil {
		return err
	}

	if err := GenerateModels(db); err != nil {
		return err
	}

	return nil
}
