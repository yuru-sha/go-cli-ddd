package repository

import (
	"context"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// ExternalAPI1AccountRepository は外部APIからアカウント情報を取得するリポジトリのインターフェースです
type ExternalAPI1AccountRepository interface {
	// FetchAccounts は外部APIから全てのアカウント情報を取得します
	FetchAccounts(ctx context.Context) ([]entity.Account, error)

	// FetchAccountByID は外部APIから指定されたIDのアカウント情報を取得します
	FetchAccountByID(ctx context.Context, id int) (entity.Account, error)
}
