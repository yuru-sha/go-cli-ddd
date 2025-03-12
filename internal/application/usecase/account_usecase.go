package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/model"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// AccountUseCase はアカウント関連のユースケースを実装します
type AccountUseCase struct {
	accountRepo      repository.MySQLAccountRepository
	accountAPIRepo   repository.ExternalAPI1AccountRepository
	notificationRepo repository.NotificationRepository
}

// NewAccountUseCase は AccountUseCase の新しいインスタンスを作成します
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

// SyncAccounts は外部APIからアカウント情報を取得し、データベースに同期します
func (uc *AccountUseCase) SyncAccounts(ctx context.Context) error {
	log.Info().Msg("アカウント情報の同期を開始します")

	// コマンド実行結果の記録を開始
	result := model.NewCommandResult("account sync")

	// 外部APIからアカウント情報を取得
	accounts, err := uc.accountAPIRepo.FetchAccounts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("アカウント情報の取得に失敗しました")
		result.SetFailed()
		result.Complete()
		if notifyErr := uc.notificationRepo.NotifyCommandResult(result); notifyErr != nil {
			log.Error().Err(notifyErr).Msg("通知の送信に失敗しました")
		}
		return err
	}

	log.Info().Int("count", len(accounts)).Msg("アカウント情報を取得しました")

	// データベースに保存
	if err := uc.accountRepo.SaveAll(ctx, accounts); err != nil {
		log.Error().Err(err).Msg("アカウント情報の保存に失敗しました")
		result.SetFailed()
		result.Complete()
		if notifyErr := uc.notificationRepo.NotifyCommandResult(result); notifyErr != nil {
			log.Error().Err(notifyErr).Msg("通知の送信に失敗しました")
		}
		return err
	}

	// 処理結果を記録
	result.AddCounts(len(accounts), 0, len(accounts))
	result.Complete()

	// 通知を送信
	if err := uc.notificationRepo.NotifyCommandResult(result); err != nil {
		log.Error().Err(err).Msg("通知の送信に失敗しました")
	}

	log.Info().Msg("アカウント情報の同期が完了しました")
	return nil
}

// SyncAccountsByIDs は指定されたアカウントIDのアカウント情報を同期します
func (uc *AccountUseCase) SyncAccountsByIDs(ctx context.Context, accountIDs []int) error {
	log.Info().Ints("account_ids", accountIDs).Msg("指定されたアカウント情報の同期を開始します")

	// コマンド実行結果の記録を開始
	result := model.NewCommandResult("account sync --id")

	// アカウントIDを文字列に変換
	accountIDStrs := make([]string, len(accountIDs))
	for i, id := range accountIDs {
		accountIDStrs[i] = strconv.Itoa(id)
	}
	result.SetAccountIDs(accountIDStrs)

	successCount := 0
	errorCount := 0
	totalRecords := 0

	// 各アカウントIDについて処理
	for _, accountID := range accountIDs {
		// 外部APIからアカウント情報を取得
		account, err := uc.accountAPIRepo.FetchAccountByID(ctx, accountID)
		if err != nil {
			log.Error().Err(err).Int("account_id", accountID).Msg("アカウント情報の取得に失敗しました")
			errorCount++
			continue
		}

		// データベースに保存
		if err := uc.accountRepo.Save(ctx, account); err != nil {
			log.Error().Err(err).Int("account_id", accountID).Msg("アカウント情報の保存に失敗しました")
			errorCount++
			continue
		}

		successCount++
		totalRecords++
	}

	// 処理結果を記録
	result.AddCounts(successCount, errorCount, totalRecords)

	if errorCount > 0 && successCount == 0 {
		result.SetFailed()
	}

	result.Complete()

	// 通知を送信
	if err := uc.notificationRepo.NotifyCommandResult(result); err != nil {
		log.Error().Err(err).Msg("通知の送信に失敗しました")
	}

	log.Info().
		Int("success", successCount).
		Int("error", errorCount).
		Int("total", len(accountIDs)).
		Msg("指定されたアカウント情報の同期が完了しました")

	if errorCount > 0 && successCount == 0 {
		return fmt.Errorf("すべてのアカウント情報の同期に失敗しました")
	}

	return nil
}

// GetAllAccounts は全てのアカウント情報を取得します
func (uc *AccountUseCase) GetAllAccounts(ctx context.Context) ([]entity.Account, error) {
	return uc.accountRepo.FindAll(ctx)
}
