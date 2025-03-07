package repository

import (
	"context"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// CampaignAPIRepository は外部APIからキャンペーン情報を取得するリポジトリのインターフェースです
type CampaignAPIRepository interface {
	// FetchCampaignsByAccountID は外部APIから指定されたアカウントIDに関連するキャンペーン情報を取得します
	FetchCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error)
}
