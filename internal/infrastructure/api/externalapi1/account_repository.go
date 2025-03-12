package externalapi1

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
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// AccountRepositoryImpl はExternalAPI1AccountRepositoryインターフェースの実装です
type AccountRepositoryImpl struct {
	client    *http.Client
	baseURL   string
	mock      bool
	apiClient *APIClient
}

// NewAccountRepository は新しいAccountRepositoryImplインスタンスを作成します
func NewAccountRepository(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) repository.ExternalAPI1AccountRepository {
	apiClient := NewAPIClient(cfg, httpClient, secretsManager)

	return &AccountRepositoryImpl{
		client:    httpClient,
		baseURL:   cfg.ExternalAPI1.BaseURL,
		mock:      true, // 常にモックを使用
		apiClient: apiClient,
	}
}

// FetchAccounts は外部APIからアカウント情報を取得します
func (r *AccountRepositoryImpl) FetchAccounts(ctx context.Context) ([]entity.Account, error) {
	if r.mock {
		return r.fetchMockAccounts(ctx)
	}

	// 実際のAPIリクエストを行う場合の実装
	url := fmt.Sprintf("%s%s", r.baseURL, "/api/accounts")

	// 認証情報を含むリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("リクエストの作成に失敗しました")
		return nil, err
	}

	// トークンを取得して設定
	headerName, headerValue, err := r.apiClient.GetAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if headerName != "" && headerValue != "" {
		req.Header.Set(headerName, headerValue)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	resp, err := r.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("APIリクエストの送信に失敗しました")
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIリクエストが失敗しました: ステータスコード %d", resp.StatusCode)
		log.Error().Err(err).Int("status_code", resp.StatusCode).Msg("APIリクエストが失敗しました")
		return nil, err
	}

	// レスポンスをデコード
	var accounts []entity.Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		log.Error().Err(err).Msg("APIレスポンスのデコードに失敗しました")
		return nil, err
	}

	return accounts, nil
}

// FetchAccountByID は外部APIから指定されたIDのアカウント情報を取得します
func (r *AccountRepositoryImpl) FetchAccountByID(ctx context.Context, id int) (entity.Account, error) {
	if r.mock {
		return r.fetchMockAccountByID(ctx, id)
	}

	// 実際のAPIリクエストを行う場合の実装
	url := fmt.Sprintf("%s%s/%d", r.baseURL, "/api/accounts", id)

	// 認証情報を含むリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("リクエストの作成に失敗しました")
		return entity.Account{}, err
	}

	// トークンを取得して設定
	headerName, headerValue, err := r.apiClient.GetAuthorizationHeader(ctx)
	if err != nil {
		return entity.Account{}, err
	}

	if headerName != "" && headerValue != "" {
		req.Header.Set(headerName, headerValue)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	resp, err := r.client.Do(req)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("APIリクエストの送信に失敗しました")
		return entity.Account{}, err
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIリクエストが失敗しました: ステータスコード %d", resp.StatusCode)
		log.Error().Err(err).Int("status_code", resp.StatusCode).Int("id", id).Msg("APIリクエストが失敗しました")
		return entity.Account{}, err
	}

	// レスポンスをデコード
	var account entity.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		log.Error().Err(err).Int("id", id).Msg("APIレスポンスのデコードに失敗しました")
		return entity.Account{}, err
	}

	return account, nil
}

// fetchMockAccounts はモックのアカウント情報を返します
func (r *AccountRepositoryImpl) fetchMockAccounts(_ context.Context) ([]entity.Account, error) {
	// モックデータを作成
	accounts := []entity.Account{
		{
			ID:        1,
			Name:      "テストアカウント1",
			Status:    "active",
			APIKey:    "api_key_test_1",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "テストアカウント2",
			Status:    "inactive",
			APIKey:    "api_key_test_2",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	log.Debug().Msg("モックアカウントデータを返します")
	return accounts, nil
}

// fetchMockAccountByID はモックの単一アカウント情報を返します
func (r *AccountRepositoryImpl) fetchMockAccountByID(_ context.Context, id int) (entity.Account, error) {
	// モックデータを作成
	accounts := map[int]entity.Account{
		1: {
			ID:        1,
			Name:      "テストアカウント1",
			Status:    "active",
			APIKey:    "api_key_test_1",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		2: {
			ID:        2,
			Name:      "テストアカウント2",
			Status:    "inactive",
			APIKey:    "api_key_test_2",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	account, exists := accounts[id]
	if !exists {
		err := fmt.Errorf("アカウントが見つかりません: ID %d", id)
		log.Error().Err(err).Int("id", id).Msg("モックアカウントが見つかりません")
		return entity.Account{}, err
	}

	log.Debug().Int("id", id).Msg("モックアカウントデータを返します")
	return account, nil
}
