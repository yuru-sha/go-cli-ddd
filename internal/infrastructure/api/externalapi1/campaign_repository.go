package externalapi1

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// CampaignRepositoryImpl はExternalAPI1CampaignRepositoryインターフェースの実装です。
type CampaignRepositoryImpl struct {
	client    *http.Client
	baseURL   string
	mock      bool
	apiClient *APIClient
}

// NewCampaignRepository creates the External API 1 campaign repository adapter.
func NewCampaignRepository(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) repository.ExternalAPI1CampaignRepository {
	apiClient := NewAPIClient(cfg, httpClient, secretsManager)

	return &CampaignRepositoryImpl{
		client:    httpClient,
		baseURL:   cfg.ExternalAPI1.BaseURL,
		mock:      true,
		apiClient: apiClient,
	}
}

// FetchCampaignsByAccountID fetches campaigns for one account from External API 1.
func (r *CampaignRepositoryImpl) FetchCampaignsByAccountID(ctx context.Context, accountID uint) ([]entity.Campaign, error) {
	if r.mock {
		return r.fetchMockCampaigns(ctx, accountID)
	}

	url := fmt.Sprintf("%s%s?account_id=%d", r.baseURL, "/api/campaigns", accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("APIリクエストの作成に失敗しました", "err", err, "url", url, "account_id", accountID)
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		slog.Error("APIリクエストの実行に失敗しました", "err", err, "url", url, "account_id", accountID)
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIがエラーを返しました: %s", resp.Status)
		slog.Error("APIエラー", "err", err, "url", url, "status_code", resp.StatusCode, "account_id", accountID)
		return nil, err
	}

	var campaigns []entity.Campaign
	if err := json.NewDecoder(resp.Body).Decode(&campaigns); err != nil {
		slog.Error("APIレスポンスのデコードに失敗しました", "err", err, "url", url, "account_id", accountID)
		return nil, err
	}

	return campaigns, nil
}

func (r *CampaignRepositoryImpl) fetchMockCampaigns(_ context.Context, accountID uint) ([]entity.Campaign, error) {
	slog.Info("モックキャンペーンデータを使用します", "account_id", accountID)

	now := time.Now()
	campaignCountBig, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		return nil, fmt.Errorf("乱数生成に失敗しました: %w", err)
	}
	campaignCount := int(campaignCountBig.Int64()) + 1

	mockCampaigns := make([]entity.Campaign, 0, campaignCount)
	statuses := []string{"active", "paused", "completed", "draft"}

	for i := 1; i <= campaignCount; i++ {
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

		startDate := now.AddDate(0, 0, -dayOffset)
		endDate := startDate.AddDate(0, 0, duration)

		var campaignID uint
		if accountID <= (^uint(0))/100 {
			campaignID = accountID * 100
		} else {
			campaignID = ^uint(0) - uint(i)
		}

		if i < 100 && i > 0 && campaignID <= (^uint(0))-uint(i) {
			campaignID += uint(i)
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

	time.Sleep(300 * time.Millisecond)

	return mockCampaigns, nil
}
