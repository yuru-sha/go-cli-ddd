package repository

import (
	"context"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// ExternalAPI1CampaignRepository は外部APIからキャンペーン情報を取得するリポジトリのインターフェースです
type ExternalAPI1CampaignRepository interface {
	// FetchCampaignsByAccountID は外部APIから指定されたアカウントIDに関連するキャンペーン情報を取得します
	FetchCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error)
}
