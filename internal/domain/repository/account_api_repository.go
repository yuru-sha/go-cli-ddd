package repository

import (
	"context"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// AccountAPIRepository は外部APIからアカウント情報を取得するリポジトリのインターフェースです
type AccountAPIRepository interface {
	// FetchAccounts は外部APIから全てのアカウント情報を取得します
	FetchAccounts(ctx context.Context) ([]entity.Account, error)
}
