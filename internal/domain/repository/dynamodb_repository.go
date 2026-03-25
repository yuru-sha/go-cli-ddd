// Package repository defines domain repository contracts.
package repository

import (
	"context"
)

// DynamoDBRepository はDynamoDBを使用したアイテム情報の永続化を担当するリポジトリのインターフェースです
type DynamoDBRepository interface {
	// GetItem は指定されたキーでアイテムを取得します
	GetItem(ctx context.Context, partitionKey string, sortKey string) (map[string]any, error)

	// PutItem は新しいアイテムを作成または更新します
	PutItem(ctx context.Context, item map[string]any) error

	// DeleteItem は指定されたキーのアイテムを削除します
	DeleteItem(ctx context.Context, partitionKey string, sortKey string) error

	// Query はパーティションキーと条件に基づいてアイテムを検索します
	Query(ctx context.Context, partitionKey string, filterExpression string) ([]map[string]any, error)

	// Scan はテーブル全体をスキャンして条件に一致するアイテムを検索します
	Scan(ctx context.Context, filterExpression string) ([]map[string]any, error)

	// BatchWrite は複数のアイテムを一括で書き込みます
	BatchWrite(ctx context.Context, items []map[string]any) error

	// TransactWrite はトランザクション内で複数の書き込み操作を実行します
	TransactWrite(ctx context.Context, operations []map[string]any) error
}
