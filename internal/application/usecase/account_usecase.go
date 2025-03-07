package usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// AccountUseCase はアカウント関連のユースケースを実装します
type AccountUseCase struct {
	accountRepo    repository.AccountRepository
	accountAPIRepo repository.AccountAPIRepository
}

// NewAccountUseCase は AccountUseCase の新しいインスタンスを作成します
func NewAccountUseCase(
	accountRepo repository.AccountRepository,
	accountAPIRepo repository.AccountAPIRepository,
) *AccountUseCase {
	return &AccountUseCase{
		accountRepo:    accountRepo,
		accountAPIRepo: accountAPIRepo,
	}
}

// SyncAccounts は外部APIからアカウント情報を取得し、データベースに同期します
func (uc *AccountUseCase) SyncAccounts(ctx context.Context) error {
	log.Info().Msg("アカウント情報の同期を開始します")

	// 外部APIからアカウント情報を取得
	accounts, err := uc.accountAPIRepo.FetchAccounts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("アカウント情報の取得に失敗しました")
		return err
	}

	log.Info().Int("count", len(accounts)).Msg("アカウント情報を取得しました")

	// データベースに保存
	if err := uc.accountRepo.SaveAll(ctx, accounts); err != nil {
		log.Error().Err(err).Msg("アカウント情報の保存に失敗しました")
		return err
	}

	log.Info().Msg("アカウント情報の同期が完了しました")
	return nil
}

// GetAllAccounts は全てのアカウント情報を取得します
func (uc *AccountUseCase) GetAllAccounts(ctx context.Context) ([]entity.Account, error) {
	return uc.accountRepo.FindAll(ctx)
}
