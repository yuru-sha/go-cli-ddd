// Package usecase implements application-level orchestration.
package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/model"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// AccountUseCase はアカウント関連のユースケースを実装します。
type AccountUseCase struct {
	accountRepo      repository.MySQLAccountRepository
	accountAPIRepo   repository.ExternalAPI1AccountRepository
	notificationRepo repository.NotificationRepository
}

// NewAccountUseCase creates an AccountUseCase.
func NewAccountUseCase(
	accountRepo repository.MySQLAccountRepository,
	accountAPIRepo repository.ExternalAPI1AccountRepository,
	notificationRepo repository.NotificationRepository,
) *AccountUseCase {
	return &AccountUseCase{
		accountRepo:      accountRepo,
		accountAPIRepo:   accountAPIRepo,
		notificationRepo: notificationRepo,
	}
}

// SyncAccounts synchronizes all accounts from the external API.
func (uc *AccountUseCase) SyncAccounts(ctx context.Context) error {
	slog.Info("アカウント情報の同期を開始します")

	result := model.NewCommandResult("account sync")

	accounts, err := uc.accountAPIRepo.FetchAccounts(ctx)
	if err != nil {
		slog.Error("アカウント情報の取得に失敗しました", "err", err)
		result.SetFailed()
		result.Complete()
		if notifyErr := uc.notificationRepo.NotifyCommandResult(result); notifyErr != nil {
			slog.Error("通知の送信に失敗しました", "err", notifyErr)
		}
		return err
	}

	slog.Info("アカウント情報を取得しました", "count", len(accounts))

	if err := uc.accountRepo.SaveAll(ctx, accounts); err != nil {
		slog.Error("アカウント情報の保存に失敗しました", "err", err)
		result.SetFailed()
		result.Complete()
		if notifyErr := uc.notificationRepo.NotifyCommandResult(result); notifyErr != nil {
			slog.Error("通知の送信に失敗しました", "err", notifyErr)
		}
		return err
	}

	result.AddCounts(len(accounts), 0, len(accounts))
	result.Complete()

	if err := uc.notificationRepo.NotifyCommandResult(result); err != nil {
		slog.Error("通知の送信に失敗しました", "err", err)
	}

	slog.Info("アカウント情報の同期が完了しました")
	return nil
}

// SyncAccountsByIDs synchronizes only the specified accounts.
func (uc *AccountUseCase) SyncAccountsByIDs(ctx context.Context, accountIDs []int) error {
	slog.Info("指定されたアカウント情報の同期を開始します", "account_ids", accountIDs)

	result := model.NewCommandResult("account sync --id")

	accountIDStrs := make([]string, len(accountIDs))
	for i, id := range accountIDs {
		accountIDStrs[i] = strconv.Itoa(id)
	}
	result.SetAccountIDs(accountIDStrs)

	successCount := 0
	errorCount := 0
	totalRecords := 0

	for _, accountID := range accountIDs {
		account, err := uc.accountAPIRepo.FetchAccountByID(ctx, accountID)
		if err != nil {
			slog.Error("アカウント情報の取得に失敗しました", "err", err, "account_id", accountID)
			errorCount++
			continue
		}

		if err := uc.accountRepo.Save(ctx, account); err != nil {
			slog.Error("アカウント情報の保存に失敗しました", "err", err, "account_id", accountID)
			errorCount++
			continue
		}

		successCount++
		totalRecords++
	}

	result.AddCounts(successCount, errorCount, totalRecords)
	if errorCount > 0 && successCount == 0 {
		result.SetFailed()
	}
	result.Complete()

	if err := uc.notificationRepo.NotifyCommandResult(result); err != nil {
		slog.Error("通知の送信に失敗しました", "err", err)
	}

	slog.Info(
		"指定されたアカウント情報の同期が完了しました",
		"success", successCount,
		"error", errorCount,
		"total", len(accountIDs),
	)

	if errorCount > 0 && successCount == 0 {
		return fmt.Errorf("すべてのアカウント情報の同期に失敗しました")
	}

	return nil
}

// GetAllAccounts fetches all stored accounts.
func (uc *AccountUseCase) GetAllAccounts(ctx context.Context) ([]entity.Account, error) {
	return uc.accountRepo.FindAll(ctx)
}
