package mysql

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// CampaignRepositoryImpl はMySQLCampaignRepositoryインターフェースの実装です。
type CampaignRepositoryImpl struct {
	db *gorm.DB
}

// NewCampaignRepository creates a MySQL-backed campaign repository.
func NewCampaignRepository(db *gorm.DB) repository.MySQLCampaignRepository {
	return &CampaignRepositoryImpl{db: db}
}

// FindAll returns all campaigns.
func (r *CampaignRepositoryImpl) FindAll(_ context.Context) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	result := r.db.Find(&campaigns)
	if result.Error != nil {
		slog.Error("キャンペーン一覧の取得に失敗しました", "err", result.Error)
		return nil, result.Error
	}
	return campaigns, nil
}

// FindByID returns one campaign by ID.
func (r *CampaignRepositoryImpl) FindByID(_ context.Context, id uint) (*entity.Campaign, error) {
	var campaign entity.Campaign
	result := r.db.First(&campaign, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			slog.Debug("キャンペーンが見つかりませんでした", "id", id)
			return nil, nil
		}
		slog.Error("キャンペーンの取得に失敗しました", "err", result.Error, "id", id)
		return nil, result.Error
	}
	return &campaign, nil
}

// FindByAccountID returns campaigns for one account.
func (r *CampaignRepositoryImpl) FindByAccountID(_ context.Context, accountID uint) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	result := r.db.Where("account_id = ?", accountID).Find(&campaigns)
	if result.Error != nil {
		slog.Error("アカウントに関連するキャンペーンの取得に失敗しました", "err", result.Error, "account_id", accountID)
		return nil, result.Error
	}
	return campaigns, nil
}

// Create inserts a new campaign.
func (r *CampaignRepositoryImpl) Create(_ context.Context, campaign *entity.Campaign) error {
	result := r.db.Create(campaign)
	if result.Error != nil {
		slog.Error("キャンペーンの作成に失敗しました", "err", result.Error)
		return result.Error
	}
	return nil
}

// Update updates an existing campaign.
func (r *CampaignRepositoryImpl) Update(_ context.Context, campaign *entity.Campaign) error {
	result := r.db.Save(campaign)
	if result.Error != nil {
		slog.Error("キャンペーンの更新に失敗しました", "err", result.Error, "id", campaign.ID)
		return result.Error
	}
	return nil
}

// Delete removes one campaign by ID.
func (r *CampaignRepositoryImpl) Delete(_ context.Context, id uint) error {
	result := r.db.Delete(&entity.Campaign{}, id)
	if result.Error != nil {
		slog.Error("キャンペーンの削除に失敗しました", "err", result.Error, "id", id)
		return result.Error
	}
	return nil
}

// SaveAll upserts a batch of campaigns in one transaction.
func (r *CampaignRepositoryImpl) SaveAll(_ context.Context, campaigns []entity.Campaign) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingCampaigns []entity.Campaign
		if err := tx.Find(&existingCampaigns).Error; err != nil {
			slog.Error("既存キャンペーンの取得に失敗しました", "err", err)
			return err
		}

		existingMap := make(map[uint]entity.Campaign)
		for _, camp := range existingCampaigns {
			existingMap[camp.ID] = camp
		}

		for _, campaign := range campaigns {
			if _, exists := existingMap[campaign.ID]; exists {
				if err := tx.Save(&campaign).Error; err != nil {
					slog.Error("キャンペーンの更新に失敗しました", "err", err, "id", campaign.ID)
					return err
				}
				continue
			}

			if err := tx.Create(&campaign).Error; err != nil {
				slog.Error("キャンペーンの作成に失敗しました", "err", err)
				return err
			}
		}

		return nil
	})
}
