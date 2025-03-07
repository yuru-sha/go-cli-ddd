package repository

import (
	"context"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// AccountRepository はアカウント情報の永続化を担当するリポジトリのインターフェースです
type AccountRepository interface {
	// FindAll は全てのアカウントを取得します
	FindAll(ctx context.Context) ([]entity.Account, error)
	
	// FindByID は指定されたIDのアカウントを取得します
	FindByID(ctx context.Context, id uint) (*entity.Account, error)
	
	// Create は新しいアカウントを作成します
	Create(ctx context.Context, account *entity.Account) error
	
	// Update は既存のアカウントを更新します
	Update(ctx context.Context, account *entity.Account) error
	
	// Delete は指定されたIDのアカウントを削除します
	Delete(ctx context.Context, id uint) error
	
	// SaveAll は複数のアカウントを一括で保存します
	SaveAll(ctx context.Context, accounts []entity.Account) error
}
