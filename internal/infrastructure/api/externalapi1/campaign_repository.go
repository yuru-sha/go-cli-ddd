package externalapi1

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// CampaignRepositoryImpl はExternalAPI1CampaignRepositoryインターフェースの実装です
type CampaignRepositoryImpl struct {
	client    *http.Client
	baseURL   string
	mock      bool
	apiClient *APIClient
}

// NewCampaignRepository は新しいCampaignRepositoryImplインスタンスを作成します
func NewCampaignRepository(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) repository.ExternalAPI1CampaignRepository {
	apiClient := NewAPIClient(cfg, httpClient, secretsManager)

	return &CampaignRepositoryImpl{
		client:    httpClient,
		baseURL:   cfg.ExternalAPI1.BaseURL,
		mock:      true, // 常にモックを使用
		apiClient: apiClient,
	}
}

// FetchCampaignsByAccountID は外部APIから指定されたアカウントIDに関連するキャンペーン情報を取得します
func (r *CampaignRepositoryImpl) FetchCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error) {
	if r.mock {
		return r.fetchMockCampaigns(ctx, accountID)
	}

	// 実際のAPIリクエストを行う場合の実装（今回は使用しない）
	url := fmt.Sprintf("%s%s?account_id=%d", r.baseURL, "/api/campaigns", accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Str("url", url).Uint("account_id", accountID).Msg("APIリクエストの作成に失敗しました")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Uint("account_id", accountID).Msg("APIリクエストの実行に失敗しました")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIがエラーを返しました: %s", resp.Status)
		log.Error().Err(err).Str("url", url).Int("status_code", resp.StatusCode).Uint("account_id", accountID).Msg("APIエラー")
		return nil, err
	}

	var campaigns []entity.Campaign
	if err := json.NewDecoder(resp.Body).Decode(&campaigns); err != nil {
		log.Error().Err(err).Str("url", url).Uint("account_id", accountID).Msg("APIレスポンスのデコードに失敗しました")
		return nil, err
	}

	return campaigns, nil
}

// fetchMockCampaigns はモックのキャンペーンデータを返します
func (r *CampaignRepositoryImpl) fetchMockCampaigns(_ context.Context, accountID uint) ([]entity.Campaign, error) {
	log.Info().Uint("account_id", accountID).Msg("モックキャンペーンデータを使用します")

	// 現在時刻
	now := time.Now()

	// 安全な乱数生成のためにcrypto/randを使用
	// キャンペーン数を決定（1〜5個）
	campaignCountBig, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
	}
	campaignCount := int(campaignCountBig.Int64()) + 1

	// モックデータ
	mockCampaigns := make([]entity.Campaign, 0, campaignCount)

	// キャンペーンステータスのリスト
	statuses := []string{"active", "paused", "completed", "draft"}

	for i := 1; i <= campaignCount; i++ {
		// 安全な乱数生成
		dayOffsetBig, err := rand.Int(rand.Reader, big.NewInt(30))
		if err != nil {
			return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
		}
		dayOffset := int(dayOffsetBig.Int64())

		durationBig, err := rand.Int(rand.Reader, big.NewInt(60))
		if err != nil {
			return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
		}
		duration := int(durationBig.Int64()) + 30

		budgetBig, err := rand.Int(rand.Reader, big.NewInt(99))
		if err != nil {
			return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
		}
		budget := 10000.0 + float64(budgetBig.Int64())*10000.0

		statusIndexBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(statuses))))
		if err != nil {
			return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
		}
		status := statuses[statusIndexBig.Int64()]

		createdOffsetBig, err := rand.Int(rand.Reader, big.NewInt(30))
		if err != nil {
			return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
		}
		createdOffset := int(createdOffsetBig.Int64())

		// 開始日と終了日を生成
		startDate := now.AddDate(0, 0, -dayOffset)
		endDate := startDate.AddDate(0, 0, duration)

		// 安全な整数変換（オーバーフロー防止）
		var campaignID uint
		if accountID <= (^uint(0))/100 { // オーバーフロー防止のチェック
			campaignID = accountID * 100
		} else {
			campaignID = ^uint(0) - uint(i) // 最大値に近い値を使用
		}

		// iが小さい値であることを確認
		if i < 100 && i > 0 { // 適当な上限と下限
			// 整数オーバーフローを防止するための追加チェック
			if campaignID <= (^uint(0))-uint(i) {
				campaignID += uint(i)
			}
		}

		campaign := entity.Campaign{
			ID:        campaignID,
			AccountID: accountID,
			Name:      fmt.Sprintf("キャンペーン%d-%d", accountID, i),
			Status:    status,
			Budget:    budget,
			StartDate: startDate,
			EndDate:   endDate,
			CreatedAt: now.Add(-time.Duration(createdOffset) * 24 * time.Hour),
			UpdatedAt: now,
		}

		mockCampaigns = append(mockCampaigns, campaign)
	}

	// 少し遅延を入れてAPIリクエストをシミュレート
	time.Sleep(300 * time.Millisecond)

	return mockCampaigns, nil
}
