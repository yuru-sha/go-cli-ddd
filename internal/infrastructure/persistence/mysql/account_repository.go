package mysql

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// AccountRepositoryImpl はMySQLAccountRepositoryインターフェースの実装です
type AccountRepositoryImpl struct {
	db *Database
}

// NewAccountRepository は新しいAccountRepositoryImplインスタンスを作成します
func NewAccountRepository(db *gorm.DB) repository.MySQLAccountRepository {
	// データベース接続をラップ
	database := &Database{
		DB: db,
	}
	return &AccountRepositoryImpl{db: database}
}

// FindAll は全てのアカウントを取得します
func (r *AccountRepositoryImpl) FindAll(_ context.Context) ([]entity.Account, error) {
	var accounts []entity.Account
	// リードレプリカを使用
	result := r.db.GetReader().Find(&accounts)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("アカウント一覧の取得に失敗しました")
		return nil, result.Error
	}
	return accounts, nil
}

// FindByID は指定されたIDのアカウントを取得します
func (r *AccountRepositoryImpl) FindByID(_ context.Context, id uint) (*entity.Account, error) {
	var account entity.Account
	// リードレプリカを使用
	result := r.db.GetReader().First(&account, id)
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
func (r *AccountRepositoryImpl) Create(_ context.Context, account *entity.Account) error {
	// ライターを使用
	result := r.db.GetWriter().Create(account)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("アカウントの作成に失敗しました")
		return result.Error
	}
	log.Debug().Uint("id", account.ID).Msg("アカウントを作成しました")
	return nil
}

// Update は既存のアカウントを更新します
func (r *AccountRepositoryImpl) Update(_ context.Context, account *entity.Account) error {
	// ライターを使用
	result := r.db.GetWriter().Save(account)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", account.ID).Msg("アカウントの更新に失敗しました")
		return result.Error
	}
	log.Debug().Uint("id", account.ID).Msg("アカウントを更新しました")
	return nil
}

// Delete は指定されたIDのアカウントを削除します
func (r *AccountRepositoryImpl) Delete(_ context.Context, id uint) error {
	// ライターを使用
	result := r.db.GetWriter().Delete(&entity.Account{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("アカウントの削除に失敗しました")
		return result.Error
	}
	log.Debug().Uint("id", id).Msg("アカウントを削除しました")
	return nil
}

// SaveAll は複数のアカウントを一括で保存します
func (r *AccountRepositoryImpl) SaveAll(_ context.Context, accounts []entity.Account) error {
	// トランザクションを開始
	tx := r.db.GetWriter().Begin()
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("トランザクションの開始に失敗しました")
		return tx.Error
	}

	// 処理が終了したらロールバックまたはコミット
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Interface("recover", r).Msg("パニックが発生したためロールバックしました")
		}
	}()

	// 各アカウントを保存
	for i := range accounts {
		if err := tx.Save(&accounts[i]).Error; err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("アカウントの保存に失敗したためロールバックしました")
			return err
		}
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("トランザクションのコミットに失敗しました")
		return err
	}

	log.Info().Int("count", len(accounts)).Msg("アカウントを一括保存しました")
	return nil
}

// Save は単一のアカウントを保存します（存在しない場合は作成、存在する場合は更新）
func (r *AccountRepositoryImpl) Save(_ context.Context, account entity.Account) error {
	// ライターを使用
	result := r.db.GetWriter().Save(&account)
	if result.Error != nil {
		log.Error().Err(result.Error).Interface("account", account).Msg("アカウントの保存に失敗しました")
		return result.Error
	}
	log.Debug().Uint("id", account.ID).Msg("アカウントを保存しました")
	return nil
}
