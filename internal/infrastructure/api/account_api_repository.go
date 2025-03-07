package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// AccountAPIRepositoryImpl はAccountAPIRepositoryインターフェースの実装です
type AccountAPIRepositoryImpl struct {
	client  *http.Client
	baseURL string
	endpoint string
	mock    bool
}

// NewAccountAPIRepository は新しいAccountAPIRepositoryImplインスタンスを作成します
func NewAccountAPIRepository(cfg *config.APIConfig, httpClient *http.Client) repository.AccountAPIRepository {
	return &AccountAPIRepositoryImpl{
		client:   httpClient,
		baseURL:  cfg.Account.BaseURL,
		endpoint: cfg.Account.Endpoint,
		mock:     true, // 常にモックを使用
	}
}

// FetchAccounts は外部APIからアカウント情報を取得します
func (r *AccountAPIRepositoryImpl) FetchAccounts(ctx context.Context) ([]entity.Account, error) {
	if r.mock {
		return r.fetchMockAccounts(ctx)
	}

	// 実際のAPIリクエストを行う場合の実装（今回は使用しない）
	url := fmt.Sprintf("%s%s", r.baseURL, r.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("APIリクエストの作成に失敗しました")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("APIリクエストの実行に失敗しました")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIがエラーを返しました: %s", resp.Status)
		log.Error().Err(err).Str("url", url).Int("status_code", resp.StatusCode).Msg("APIエラー")
		return nil, err
	}

	var accounts []entity.Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		log.Error().Err(err).Str("url", url).Msg("APIレスポンスのデコードに失敗しました")
		return nil, err
	}

	return accounts, nil
}

// fetchMockAccounts はモックのアカウントデータを返します
func (r *AccountAPIRepositoryImpl) fetchMockAccounts(ctx context.Context) ([]entity.Account, error) {
	log.Info().Msg("モックアカウントデータを使用します")
	
	// 現在時刻
	now := time.Now()
	
	// モックデータ
	mockAccounts := []entity.Account{
		{
			ID:        1,
			Name:      "株式会社サンプル",
			Status:    "active",
			APIKey:    "api_key_sample_1",
			CreatedAt: now.Add(-30 * 24 * time.Hour),
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "テスト株式会社",
			Status:    "active",
			APIKey:    "api_key_sample_2",
			CreatedAt: now.Add(-20 * 24 * time.Hour),
			UpdatedAt: now,
		},
		{
			ID:        3,
			Name:      "広告主A",
			Status:    "inactive",
			APIKey:    "api_key_sample_3",
			CreatedAt: now.Add(-10 * 24 * time.Hour),
			UpdatedAt: now,
		},
		{
			ID:        4,
			Name:      "広告主B",
			Status:    "active",
			APIKey:    "api_key_sample_4",
			CreatedAt: now.Add(-5 * 24 * time.Hour),
			UpdatedAt: now,
		},
		{
			ID:        5,
			Name:      "広告主C",
			Status:    "pending",
			APIKey:    "api_key_sample_5",
			CreatedAt: now.Add(-1 * 24 * time.Hour),
			UpdatedAt: now,
		},
	}
	
	// 少し遅延を入れてAPIリクエストをシミュレート
	time.Sleep(500 * time.Millisecond)
	
	return mockAccounts, nil
}
