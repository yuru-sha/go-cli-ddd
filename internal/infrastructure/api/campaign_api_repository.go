package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// CampaignAPIRepositoryImpl はCampaignAPIRepositoryインターフェースの実装です
type CampaignAPIRepositoryImpl struct {
	client   *http.Client
	baseURL  string
	endpoint string
	mock     bool
}

// NewCampaignAPIRepository は新しいCampaignAPIRepositoryImplインスタンスを作成します
func NewCampaignAPIRepository(cfg *config.APIConfig, httpClient *http.Client) repository.CampaignAPIRepository {
	return &CampaignAPIRepositoryImpl{
		client:   httpClient,
		baseURL:  cfg.Campaign.BaseURL,
		endpoint: cfg.Campaign.Endpoint,
		mock:     true, // 常にモックを使用
	}
}

// FetchCampaignsByAccountID は外部APIから指定されたアカウントIDに関連するキャンペーン情報を取得します
func (r *CampaignAPIRepositoryImpl) FetchCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error) {
	if r.mock {
		return r.fetchMockCampaigns(ctx, accountID)
	}

	// 実際のAPIリクエストを行う場合の実装（今回は使用しない）
	url := fmt.Sprintf("%s%s?account_id=%d", r.baseURL, r.endpoint, accountID)
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
func (r *CampaignAPIRepositoryImpl) fetchMockCampaigns(ctx context.Context, accountID uint) ([]entity.Campaign, error) {
	log.Info().Uint("account_id", accountID).Msg("モックキャンペーンデータを使用します")
	
	// 現在時刻
	now := time.Now()
	
	// 乱数生成器の初期化
	rnd := rand.New(rand.NewSource(int64(accountID)))
	
	// キャンペーン数を決定（1〜5個）
	campaignCount := rnd.Intn(5) + 1
	
	// モックデータ
	mockCampaigns := make([]entity.Campaign, 0, campaignCount)
	
	// キャンペーンステータスのリスト
	statuses := []string{"active", "paused", "completed", "draft"}
	
	for i := 1; i <= campaignCount; i++ {
		// 開始日と終了日を生成
		startDate := now.AddDate(0, 0, -rnd.Intn(30))
		endDate := startDate.AddDate(0, 0, 30+rnd.Intn(60))
		
		// 予算を生成（10000〜1000000円）
		budget := 10000.0 + float64(rnd.Intn(99))*10000.0
		
		// ステータスをランダムに選択
		status := statuses[rnd.Intn(len(statuses))]
		
		campaign := entity.Campaign{
			ID:        uint(accountID*100 + uint(i)),
			AccountID: accountID,
			Name:      fmt.Sprintf("キャンペーン%d-%d", accountID, i),
			Status:    status,
			Budget:    budget,
			StartDate: startDate,
			EndDate:   endDate,
			CreatedAt: now.Add(-time.Duration(rnd.Intn(30)) * 24 * time.Hour),
			UpdatedAt: now,
		}
		
		mockCampaigns = append(mockCampaigns, campaign)
	}
	
	// 少し遅延を入れてAPIリクエストをシミュレート
	time.Sleep(300 * time.Millisecond)
	
	return mockCampaigns, nil
}
