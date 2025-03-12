package mysql

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// CampaignRepositoryImpl はMySQLCampaignRepositoryインターフェースの実装です
type CampaignRepositoryImpl struct {
	db *gorm.DB
}

// NewCampaignRepository は新しいCampaignRepositoryImplインスタンスを作成します
func NewCampaignRepository(db *gorm.DB) repository.MySQLCampaignRepository {
	return &CampaignRepositoryImpl{db: db}
}

// FindAll は全てのキャンペーンを取得します
func (r *CampaignRepositoryImpl) FindAll(_ context.Context) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	result := r.db.Find(&campaigns)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("キャンペーン一覧の取得に失敗しました")
		return nil, result.Error
	}
	return campaigns, nil
}

// FindByID は指定されたIDのキャンペーンを取得します
func (r *CampaignRepositoryImpl) FindByID(_ context.Context, id uint) (*entity.Campaign, error) {
	var campaign entity.Campaign
	result := r.db.First(&campaign, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().Uint("id", id).Msg("キャンペーンが見つかりませんでした")
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("キャンペーンの取得に失敗しました")
		return nil, result.Error
	}
	return &campaign, nil
}

// FindByAccountID は指定されたアカウントIDに関連するキャンペーンを全て取得します
func (r *CampaignRepositoryImpl) FindByAccountID(_ context.Context, accountID uint) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	result := r.db.Where("account_id = ?", accountID).Find(&campaigns)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("account_id", accountID).Msg("アカウントに関連するキャンペーンの取得に失敗しました")
		return nil, result.Error
	}
	return campaigns, nil
}

// Create は新しいキャンペーンを作成します
func (r *CampaignRepositoryImpl) Create(_ context.Context, campaign *entity.Campaign) error {
	result := r.db.Create(campaign)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("キャンペーンの作成に失敗しました")
		return result.Error
	}
	return nil
}

// Update は既存のキャンペーンを更新します
func (r *CampaignRepositoryImpl) Update(_ context.Context, campaign *entity.Campaign) error {
	result := r.db.Save(campaign)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", campaign.ID).Msg("キャンペーンの更新に失敗しました")
		return result.Error
	}
	return nil
}

// Delete は指定されたIDのキャンペーンを削除します
func (r *CampaignRepositoryImpl) Delete(_ context.Context, id uint) error {
	result := r.db.Delete(&entity.Campaign{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("キャンペーンの削除に失敗しました")
		return result.Error
	}
	return nil
}

// SaveAll は複数のキャンペーンを一括で保存します
func (r *CampaignRepositoryImpl) SaveAll(_ context.Context, campaigns []entity.Campaign) error {
	// トランザクションを開始
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 既存のデータを取得
		var existingCampaigns []entity.Campaign
		if err := tx.Find(&existingCampaigns).Error; err != nil {
			log.Error().Err(err).Msg("既存キャンペーンの取得に失敗しました")
			return err
		}

		// 既存のIDをマップに格納
		existingMap := make(map[uint]entity.Campaign)
		for _, camp := range existingCampaigns {
			existingMap[camp.ID] = camp
		}

		// 新しいキャンペーンを作成または更新
		for _, campaign := range campaigns {
			if _, exists := existingMap[campaign.ID]; exists {
				// 既存のキャンペーンを更新
				if err := tx.Save(&campaign).Error; err != nil {
					log.Error().Err(err).Uint("id", campaign.ID).Msg("キャンペーンの更新に失敗しました")
					return err
				}
			} else {
				// 新しいキャンペーンを作成
				if err := tx.Create(&campaign).Error; err != nil {
					log.Error().Err(err).Msg("キャンペーンの作成に失敗しました")
					return err
				}
			}
		}

		return nil
	})
}
