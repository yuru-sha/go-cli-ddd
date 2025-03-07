package usecase

import (
	"context"

	"github.com/rs/zerolog/log"
)

// MasterUseCase はマスター同期関連のユースケースを実装します
type MasterUseCase struct {
	accountUseCase  *AccountUseCase
	campaignUseCase *CampaignUseCase
}

// NewMasterUseCase は MasterUseCase の新しいインスタンスを作成します
func NewMasterUseCase(
	accountUseCase *AccountUseCase,
	campaignUseCase *CampaignUseCase,
) *MasterUseCase {
	return &MasterUseCase{
		accountUseCase:  accountUseCase,
		campaignUseCase: campaignUseCase,
	}
}

// SyncAll はアカウント情報とキャンペーン情報を順番に同期します
func (uc *MasterUseCase) SyncAll(ctx context.Context) error {
	log.Info().Msg("マスター同期を開始します")

	// アカウント情報の同期
	if err := uc.accountUseCase.SyncAccounts(ctx); err != nil {
		log.Error().Err(err).Msg("アカウント情報の同期に失敗しました")
		return err
	}

	// キャンペーン情報の同期
	if err := uc.campaignUseCase.SyncCampaigns(ctx); err != nil {
		log.Error().Err(err).Msg("キャンペーン情報の同期に失敗しました")
		return err
	}

	log.Info().Msg("マスター同期が完了しました")
	return nil
}
