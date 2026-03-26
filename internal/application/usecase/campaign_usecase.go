package usecase

import (
	"context"
	"log/slog"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// CampaignUseCase はキャンペーン関連のユースケースを実装します。
type CampaignUseCase struct {
	campaignRepo    repository.MySQLCampaignRepository
	campaignAPIRepo repository.ExternalAPI1CampaignRepository
	accountRepo     repository.MySQLAccountRepository
}

// NewCampaignUseCase creates a CampaignUseCase.
func NewCampaignUseCase(
	campaignRepo repository.MySQLCampaignRepository,
	campaignAPIRepo repository.ExternalAPI1CampaignRepository,
	accountRepo repository.MySQLAccountRepository,
) *CampaignUseCase {
	return &CampaignUseCase{
		campaignRepo:    campaignRepo,
		campaignAPIRepo: campaignAPIRepo,
		accountRepo:     accountRepo,
	}
}

// SyncCampaigns synchronizes campaigns for all known accounts.
func (uc *CampaignUseCase) SyncCampaigns(ctx context.Context) error {
	slog.Info("キャンペーン情報の同期を開始します")

	accounts, err := uc.accountRepo.FindAll(ctx)
	if err != nil {
		slog.Error("アカウント情報の取得に失敗しました", "err", err)
		return err
	}

	slog.Info("アカウント情報を取得しました", "account_count", len(accounts))

	g, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	allCampaigns := make([]entity.Campaign, 0)

	for _, account := range accounts {
		g.Go(func(account entity.Account) func() error {
			return func() error {
				campaigns, err := uc.campaignAPIRepo.FetchCampaignsByAccountID(ctx, account.ID)
				if err != nil {
					slog.Error("キャンペーン情報の取得に失敗しました", "err", err, "account_id", account.ID)
					return err
				}

				slog.Info("キャンペーン情報を取得しました", "account_id", account.ID, "campaign_count", len(campaigns))

				mu.Lock()
				allCampaigns = append(allCampaigns, campaigns...)
				mu.Unlock()

				return nil
			}
		}(account))
	}

	if err := g.Wait(); err != nil {
		slog.Error("キャンペーン情報の同期中にエラーが発生しました", "err", err)
		return err
	}

	if len(allCampaigns) > 0 {
		if err := uc.campaignRepo.SaveAll(ctx, allCampaigns); err != nil {
			slog.Error("キャンペーン情報の保存に失敗しました", "err", err)
			return err
		}
	}

	slog.Info("キャンペーン情報の同期が完了しました", "total_campaigns", len(allCampaigns))
	return nil
}

// GetCampaignsByAccountID fetches stored campaigns for one account.
func (uc *CampaignUseCase) GetCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error) {
	return uc.campaignRepo.FindByAccountID(ctx, accountID)
}
