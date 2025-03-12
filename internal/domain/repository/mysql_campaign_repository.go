package repository

import (
	"context"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// MySQLCampaignRepository はキャンペーン情報の永続化を担当するリポジトリのインターフェースです
type MySQLCampaignRepository interface {
	// FindAll は全てのキャンペーンを取得します
	FindAll(ctx context.Context) ([]entity.Campaign, error)

	// FindByID は指定されたIDのキャンペーンを取得します
	FindByID(ctx context.Context, id uint) (*entity.Campaign, error)

	// FindByAccountID は指定されたアカウントIDに関連するキャンペーンを全て取得します
	FindByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error)

	// Create は新しいキャンペーンを作成します
	Create(ctx context.Context, campaign *entity.Campaign) error

	// Update は既存のキャンペーンを更新します
	Update(ctx context.Context, campaign *entity.Campaign) error

	// Delete は指定されたIDのキャンペーンを削除します
	Delete(ctx context.Context, id uint) error

	// SaveAll は複数のキャンペーンを一括で保存します
	SaveAll(ctx context.Context, campaigns []entity.Campaign) error
}
