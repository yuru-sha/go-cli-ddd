package usecase

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// CampaignUseCase はキャンペーン関連のユースケースを実装します
type CampaignUseCase struct {
	campaignRepo    repository.MySQLCampaignRepository
	campaignAPIRepo repository.ExternalAPI1CampaignRepository
	accountRepo     repository.MySQLAccountRepository
}

// NewCampaignUseCase は CampaignUseCase の新しいインスタンスを作成します
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

// SyncCampaigns は全てのアカウントに対して、外部APIからキャンペーン情報を取得し、データベースに同期します
func (uc *CampaignUseCase) SyncCampaigns(ctx context.Context) error {
	log.Info().Msg("キャンペーン情報の同期を開始します")

	// アカウント情報を取得
	accounts, err := uc.accountRepo.FindAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("アカウント情報の取得に失敗しました")
		return err
	}

	log.Info().Int("account_count", len(accounts)).Msg("アカウント情報を取得しました")

	// 並列処理のためのエラーグループを作成
	g, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	allCampaigns := make([]entity.Campaign, 0)

	// 各アカウントに対して並列処理
	for _, account := range accounts {
		account := account // ゴルーチン内で使用するためにローカル変数にコピー
		g.Go(func() error {
			// 外部APIからキャンペーン情報を取得
			campaigns, err := uc.campaignAPIRepo.FetchCampaignsByAccountID(ctx, account.ID)
			if err != nil {
				log.Error().Err(err).Uint("account_id", account.ID).Msg("キャンペーン情報の取得に失敗しました")
				return err
			}

			log.Info().Uint("account_id", account.ID).Int("campaign_count", len(campaigns)).Msg("キャンペーン情報を取得しました")

			// 結果をマージ
			mu.Lock()
			allCampaigns = append(allCampaigns, campaigns...)
			mu.Unlock()

			return nil
		})
	}

	// 全ての並列処理が完了するのを待つ
	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("キャンペーン情報の同期中にエラーが発生しました")
		return err
	}

	// データベースに保存
	if len(allCampaigns) > 0 {
		if err := uc.campaignRepo.SaveAll(ctx, allCampaigns); err != nil {
			log.Error().Err(err).Msg("キャンペーン情報の保存に失敗しました")
			return err
		}
	}

	log.Info().Int("total_campaigns", len(allCampaigns)).Msg("キャンペーン情報の同期が完了しました")
	return nil
}

// GetCampaignsByAccountID は指定されたアカウントIDに関連するキャンペーン情報を取得します
func (uc *CampaignUseCase) GetCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error) {
	return uc.campaignRepo.FindByAccountID(ctx, accountID)
}
