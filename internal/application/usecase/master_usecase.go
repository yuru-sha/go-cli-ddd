package usecase

import (
	"context"
	"log/slog"
)

// MasterUseCase はマスター同期関連のユースケースを実装します。
type MasterUseCase struct {
	accountUseCase  *AccountUseCase
	campaignUseCase *CampaignUseCase
}

// NewMasterUseCase creates a MasterUseCase.
func NewMasterUseCase(
	accountUseCase *AccountUseCase,
	campaignUseCase *CampaignUseCase,
) *MasterUseCase {
	return &MasterUseCase{
		accountUseCase:  accountUseCase,
		campaignUseCase: campaignUseCase,
	}
}

// SyncAll synchronizes all master data in sequence.
func (uc *MasterUseCase) SyncAll(ctx context.Context) error {
	slog.Info("マスター同期を開始します")

	if err := uc.accountUseCase.SyncAccounts(ctx); err != nil {
		slog.Error("アカウント情報の同期に失敗しました", "err", err)
		return err
	}

	if err := uc.campaignUseCase.SyncCampaigns(ctx); err != nil {
		slog.Error("キャンペーン情報の同期に失敗しました", "err", err)
		return err
	}

	slog.Info("マスター同期が完了しました")
	return nil
}
