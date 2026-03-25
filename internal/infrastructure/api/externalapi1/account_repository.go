package externalapi1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// AccountRepositoryImpl はExternalAPI1AccountRepositoryインターフェースの実装です。
type AccountRepositoryImpl struct {
	client    *http.Client
	baseURL   string
	mock      bool
	apiClient *APIClient
}

// NewAccountRepository creates the External API 1 account repository adapter.
func NewAccountRepository(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) repository.ExternalAPI1AccountRepository {
	apiClient := NewAPIClient(cfg, httpClient, secretsManager)

	return &AccountRepositoryImpl{
		client:    httpClient,
		baseURL:   cfg.ExternalAPI1.BaseURL,
		mock:      true,
		apiClient: apiClient,
	}
}

// FetchAccounts fetches all accounts from External API 1.
func (r *AccountRepositoryImpl) FetchAccounts(ctx context.Context) ([]entity.Account, error) {
	if r.mock {
		return r.fetchMockAccounts(ctx)
	}

	url := fmt.Sprintf("%s%s", r.baseURL, "/api/accounts")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("リクエストの作成に失敗しました", "err", err)
		return nil, err
	}

	headerName, headerValue, err := r.apiClient.GetAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if headerName != "" && headerValue != "" {
		req.Header.Set(headerName, headerValue)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		slog.Error("APIリクエストの送信に失敗しました", "err", err)
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIリクエストが失敗しました: ステータスコード %d", resp.StatusCode)
		slog.Error("APIリクエストが失敗しました", "err", err, "status_code", resp.StatusCode)
		return nil, err
	}

	var accounts []entity.Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		slog.Error("APIレスポンスのデコードに失敗しました", "err", err)
		return nil, err
	}

	return accounts, nil
}

// FetchAccountByID fetches one account by ID from External API 1.
func (r *AccountRepositoryImpl) FetchAccountByID(ctx context.Context, id int) (entity.Account, error) {
	if r.mock {
		return r.fetchMockAccountByID(ctx, id)
	}

	url := fmt.Sprintf("%s%s/%d", r.baseURL, "/api/accounts", id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("リクエストの作成に失敗しました", "err", err, "id", id)
		return entity.Account{}, err
	}

	headerName, headerValue, err := r.apiClient.GetAuthorizationHeader(ctx)
	if err != nil {
		return entity.Account{}, err
	}

	if headerName != "" && headerValue != "" {
		req.Header.Set(headerName, headerValue)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		slog.Error("APIリクエストの送信に失敗しました", "err", err, "id", id)
		return entity.Account{}, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("APIリクエストが失敗しました: ステータスコード %d", resp.StatusCode)
		slog.Error("APIリクエストが失敗しました", "err", err, "status_code", resp.StatusCode, "id", id)
		return entity.Account{}, err
	}

	var account entity.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		slog.Error("APIレスポンスのデコードに失敗しました", "err", err, "id", id)
		return entity.Account{}, err
	}

	return account, nil
}

func (r *AccountRepositoryImpl) fetchMockAccounts(_ context.Context) ([]entity.Account, error) {
	accounts := []entity.Account{
		{
			ID:        1,
			Name:      "テストアカウント1",
			Status:    "active",
			APIKey:    "mock-account-value-1",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "テストアカウント2",
			Status:    "inactive",
			APIKey:    "mock-account-value-2",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	slog.Debug("モックアカウントデータを返します")
	return accounts, nil
}

func (r *AccountRepositoryImpl) fetchMockAccountByID(_ context.Context, id int) (entity.Account, error) {
	accounts := map[int]entity.Account{
		1: {
			ID:        1,
			Name:      "テストアカウント1",
			Status:    "active",
			APIKey:    "mock-account-value-1",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		2: {
			ID:        2,
			Name:      "テストアカウント2",
			Status:    "inactive",
			APIKey:    "mock-account-value-2",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	account, exists := accounts[id]
	if !exists {
		err := fmt.Errorf("アカウントが見つかりません: ID %d", id)
		slog.Error("モックアカウントが見つかりません", "err", err, "id", id)
		return entity.Account{}, err
	}

	slog.Debug("モックアカウントデータを返します", "id", id)
	return account, nil
}
