//go:build integration
// +build integration

package mysql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// TestDatabaseConnection は統合テストの例です
// このテストは `go test -tags=integration` コマンドでのみ実行されます
func TestDatabaseConnection(t *testing.T) {
	// テスト用の設定
	dbConfig := &config.DatabaseConfig{
		Dialect:     "sqlite3",
		DSN:         ":memory:", // インメモリデータベースを使用
		LogLevel:    "silent",
		AutoMigrate: true,
	}

	// データベース接続
	db, err := NewDatabase(dbConfig)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// リポジトリの作成
	repo := NewAccountRepository(db)
	assert.NotNil(t, repo)

	// テスト用のアカウントデータ
	testAccount := &entity.Account{
		Name:   "テスト株式会社",
		Status: "active",
		APIKey: "test_api_key",
	}

	// アカウントの作成
	err = repo.Create(context.Background(), testAccount)
	assert.NoError(t, err)
	assert.NotZero(t, testAccount.ID)

	// アカウントの取得
	retrievedAccount, err := repo.FindByID(context.Background(), testAccount.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedAccount)
	assert.Equal(t, testAccount.Name, retrievedAccount.Name)
	assert.Equal(t, testAccount.Status, retrievedAccount.Status)
	assert.Equal(t, testAccount.APIKey, retrievedAccount.APIKey)
}
