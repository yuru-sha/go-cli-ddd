package persistence

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// AccountRepositoryImpl はAccountRepositoryインターフェースの実装です
type AccountRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountRepository は新しいAccountRepositoryImplインスタンスを作成します
func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &AccountRepositoryImpl{db: db}
}

// FindAll は全てのアカウントを取得します
func (r *AccountRepositoryImpl) FindAll(ctx context.Context) ([]entity.Account, error) {
	var accounts []entity.Account
	result := r.db.Find(&accounts)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("アカウント一覧の取得に失敗しました")
		return nil, result.Error
	}
	return accounts, nil
}

// FindByID は指定されたIDのアカウントを取得します
func (r *AccountRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Account, error) {
	var account entity.Account
	result := r.db.First(&account, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().Uint("id", id).Msg("アカウントが見つかりませんでした")
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("アカウントの取得に失敗しました")
		return nil, result.Error
	}
	return &account, nil
}

// Create は新しいアカウントを作成します
func (r *AccountRepositoryImpl) Create(ctx context.Context, account *entity.Account) error {
	result := r.db.Create(account)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("アカウントの作成に失敗しました")
		return result.Error
	}
	return nil
}

// Update は既存のアカウントを更新します
func (r *AccountRepositoryImpl) Update(ctx context.Context, account *entity.Account) error {
	result := r.db.Save(account)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", account.ID).Msg("アカウントの更新に失敗しました")
		return result.Error
	}
	return nil
}

// Delete は指定されたIDのアカウントを削除します
func (r *AccountRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.Account{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("アカウントの削除に失敗しました")
		return result.Error
	}
	return nil
}

// SaveAll は複数のアカウントを一括で保存します
func (r *AccountRepositoryImpl) SaveAll(ctx context.Context, accounts []entity.Account) error {
	// トランザクションを開始
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 既存のデータを取得
		var existingAccounts []entity.Account
		if err := tx.Find(&existingAccounts).Error; err != nil {
			log.Error().Err(err).Msg("既存アカウントの取得に失敗しました")
			return err
		}

		// 既存のIDをマップに格納
		existingMap := make(map[uint]entity.Account)
		for _, acc := range existingAccounts {
			existingMap[acc.ID] = acc
		}

		// 新しいアカウントを作成または更新
		for _, account := range accounts {
			if _, exists := existingMap[account.ID]; exists {
				// 既存のアカウントを更新
				if err := tx.Save(&account).Error; err != nil {
					log.Error().Err(err).Uint("id", account.ID).Msg("アカウントの更新に失敗しました")
					return err
				}
			} else {
				// 新しいアカウントを作成
				if err := tx.Create(&account).Error; err != nil {
					log.Error().Err(err).Msg("アカウントの作成に失敗しました")
					return err
				}
			}
		}

		return nil
	})
}
