package mysql

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// AccountRepositoryImpl はMySQLAccountRepositoryインターフェースの実装です。
type AccountRepositoryImpl struct {
	db *Database
}

// NewAccountRepository creates a MySQL-backed account repository.
func NewAccountRepository(db *gorm.DB) repository.MySQLAccountRepository {
	database := &Database{DB: db}
	return &AccountRepositoryImpl{db: database}
}

// FindAll returns all accounts.
func (r *AccountRepositoryImpl) FindAll(_ context.Context) ([]entity.Account, error) {
	var accounts []entity.Account
	result := r.db.GetReader().Find(&accounts)
	if result.Error != nil {
		slog.Error("アカウント一覧の取得に失敗しました", "err", result.Error)
		return nil, result.Error
	}
	return accounts, nil
}

// FindByID returns one account by ID.
func (r *AccountRepositoryImpl) FindByID(_ context.Context, id uint) (*entity.Account, error) {
	var account entity.Account
	result := r.db.GetReader().First(&account, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			slog.Debug("アカウントが見つかりませんでした", "id", id)
			return nil, nil
		}
		slog.Error("アカウントの取得に失敗しました", "err", result.Error, "id", id)
		return nil, result.Error
	}
	return &account, nil
}

// Create inserts a new account.
func (r *AccountRepositoryImpl) Create(_ context.Context, account *entity.Account) error {
	result := r.db.GetWriter().Create(account)
	if result.Error != nil {
		slog.Error("アカウントの作成に失敗しました", "err", result.Error)
		return result.Error
	}
	slog.Debug("アカウントを作成しました", "id", account.ID)
	return nil
}

// Update updates an existing account.
func (r *AccountRepositoryImpl) Update(_ context.Context, account *entity.Account) error {
	result := r.db.GetWriter().Save(account)
	if result.Error != nil {
		slog.Error("アカウントの更新に失敗しました", "err", result.Error, "id", account.ID)
		return result.Error
	}
	slog.Debug("アカウントを更新しました", "id", account.ID)
	return nil
}

// Delete removes one account by ID.
func (r *AccountRepositoryImpl) Delete(_ context.Context, id uint) error {
	result := r.db.GetWriter().Delete(&entity.Account{}, id)
	if result.Error != nil {
		slog.Error("アカウントの削除に失敗しました", "err", result.Error, "id", id)
		return result.Error
	}
	slog.Debug("アカウントを削除しました", "id", id)
	return nil
}

// SaveAll upserts a batch of accounts in one transaction.
func (r *AccountRepositoryImpl) SaveAll(_ context.Context, accounts []entity.Account) error {
	tx := r.db.GetWriter().Begin()
	if tx.Error != nil {
		slog.Error("トランザクションの開始に失敗しました", "err", tx.Error)
		return tx.Error
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			tx.Rollback()
			slog.Error("パニックが発生したためロールバックしました", "recover", recovered)
		}
	}()

	for i := range accounts {
		if err := tx.Save(&accounts[i]).Error; err != nil {
			tx.Rollback()
			slog.Error("アカウントの保存に失敗したためロールバックしました", "err", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error("トランザクションのコミットに失敗しました", "err", err)
		return err
	}

	slog.Info("アカウントを一括保存しました", "count", len(accounts))
	return nil
}

// Save upserts one account.
func (r *AccountRepositoryImpl) Save(_ context.Context, account entity.Account) error {
	result := r.db.GetWriter().Save(&account)
	if result.Error != nil {
		slog.Error("アカウントの保存に失敗しました", "err", result.Error, "account", account)
		return result.Error
	}
	slog.Debug("アカウントを保存しました", "id", account.ID)
	return nil
}
