package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// メインエントリーポイント
func main() {
	// コマンドライン引数の解析
	var host string
	var port int
	var user string
	var password string
	var dbName string
	var outPath string
	var modelPath string
	var tables string

	flag.StringVar(&host, "host", "localhost", "MySQLホスト")
	flag.IntVar(&port, "port", 3306, "MySQLポート")
	flag.StringVar(&user, "user", "root", "MySQLユーザー名")
	flag.StringVar(&password, "password", "", "MySQLパスワード")
	flag.StringVar(&dbName, "dbname", "", "データベース名（必須）")
	flag.StringVar(&outPath, "out", "./internal/infrastructure/persistence/query", "生成されるクエリコードの出力先")
	flag.StringVar(&modelPath, "model", "./internal/infrastructure/persistence/model", "生成されるモデルコードの出力先")
	flag.StringVar(&tables, "tables", "", "生成対象のテーブル名（カンマ区切り、空の場合は全テーブル）")
	flag.Parse()

	// データベース名は必須
	if dbName == "" {
		log.Error().Msg("データベース名を指定してください（-dbname）")
		flag.Usage()
		os.Exit(1)
	}

	log.Info().Msg("GORM genによるモデル生成を開始します")
	log.Info().Str("database", dbName).Str("host", host).Int("port", port).Msg("対象データベース")

	// データベース接続
	db, err := connectDatabase(host, port, user, password, dbName)
	if err != nil {
		log.Error().Err(err).Msg("データベース接続に失敗しました")
		os.Exit(1)
	}

	// モデル生成
	if err := generateModels(db, outPath, modelPath, tables); err != nil {
		log.Error().Err(err).Msg("モデル生成に失敗しました")
		os.Exit(1)
	}

	log.Info().Msg("モデル生成が完了しました")
}

// connectDatabase はデータベースに接続します
func connectDatabase(host string, port int, user, password, dbName string) (*gorm.DB, error) {
	// MySQL接続設定
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // テーブル名を単数形にする
		},
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続に失敗しました: %w", err)
	}

	return db, nil
}

// generateModels はGORM genを使用してモデルを生成します
func generateModels(db *gorm.DB, outPath, modelPath, tableList string) error {
	// ジェネレーターの設定
	g := gen.NewGenerator(gen.Config{
		OutPath:           outPath,
		ModelPkgPath:      modelPath,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true, // NULL可能なフィールドにポインタを使用
		FieldCoverable:    true, // フィールドの上書きを許可
		FieldSignable:     true, // 符号付き整数を使用
		FieldWithIndexTag: true, // インデックスタグを生成
		FieldWithTypeTag:  true, // タイプタグを生成
	})

	// データベース接続の設定
	g.UseDB(db)

	// 対象テーブルの取得
	var tables []string
	if tableList != "" {
		tables = strings.Split(tableList, ",")
		log.Info().Strs("tables", tables).Msg("指定されたテーブルのみを生成します")
	} else {
		// データベースから全テーブルを取得
		var tableNames []string
		if err := db.Raw("SHOW TABLES").Scan(&tableNames).Error; err != nil {
			return fmt.Errorf("テーブル一覧の取得に失敗しました: %w", err)
		}
		tables = tableNames
		log.Info().Strs("tables", tables).Msg("データベース内の全テーブルを生成します")
	}

	// 各テーブルに対してモデルを生成
	for _, tableName := range tables {
		// 日時型フィールドの特別処理
		timeFields := []string{"created_at", "updated_at", "start_date", "end_date", "deleted_at"}
		fieldOpts := []gen.ModelOpt{}

		for _, field := range timeFields {
			fieldOpts = append(fieldOpts, gen.FieldType(field, "time.Time"))
		}

		g.GenerateModel(tableName, fieldOpts...)
	}

	// コードの生成
	g.Execute()

	return nil
}
