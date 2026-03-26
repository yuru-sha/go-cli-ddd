package externalapi1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// APIClient は外部API1のクライアントです。
type APIClient struct {
	client         *http.Client
	config         *config.Config
	secretsManager secrets.Manager
}

// TokenCache はトークンをキャッシュするための構造体です。
type TokenCache struct {
	Token     string
	ExpiresAt time.Time
}

// NewAPIClient creates an External API 1 client.
func NewAPIClient(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) *APIClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &APIClient{
		client:         httpClient,
		config:         cfg,
		secretsManager: secretsManager,
	}
}

// GetAuthorizationHeader resolves the authorization header for External API 1.
func (c *APIClient) GetAuthorizationHeader(ctx context.Context) (string, string, error) {
	if c.config.UseAWSSecretsManager() && c.config.ExternalAPI1.TokenSecretID != "" {
		slog.Info("Secret ManagerからAPIトークンを取得します")

		tokenStr, err := c.secretsManager.GetSecret(ctx, c.config.ExternalAPI1.TokenSecretID)
		if err != nil {
			return "", "", fmt.Errorf("APIトークンの取得に失敗しました: %w", err)
		}

		var tokenSecret secrets.APITokenSecret
		if err := json.Unmarshal([]byte(tokenStr), &tokenSecret); err != nil {
			return "", "", fmt.Errorf("APIトークンのパースに失敗しました: %w", err)
		}

		if tokenSecret.BearerToken != "" {
			return "Authorization", fmt.Sprintf("Bearer %s", tokenSecret.BearerToken), nil
		}
		if tokenSecret.Token != "" {
			return "X-API-Token", tokenSecret.Token, nil
		}
		if tokenSecret.AccessKey != "" && tokenSecret.SecretKey != "" {
			return "X-API-Key", tokenSecret.AccessKey, nil
		}
	}

	if c.config.ExternalAPI1.Token != "" {
		return "X-API-Token", c.config.ExternalAPI1.Token, nil
	}

	return "X-API-Token", "dev-token-12345", nil
}

// CreateAuthenticatedRequest builds an authenticated request for External API 1.
func (c *APIClient) CreateAuthenticatedRequest(ctx context.Context, method, url string, _ any) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	headerName, headerValue, err := c.GetAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if headerName != "" && headerValue != "" {
		req.Header.Set(headerName, headerValue)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Request executes a request against External API 1.
func (c *APIClient) Request(ctx context.Context, method, path string, _ io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.config.ExternalAPI1.BaseURL, path)

	req, err := c.CreateAuthenticatedRequest(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの実行に失敗しました: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close() //nolint:errcheck
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: %s - %s", resp.Status, string(body))
	}

	return resp, nil
}

// GetData fetches one sample resource from External API 1.
func (c *APIClient) GetData(ctx context.Context, dataID string) (map[string]any, error) {
	path := fmt.Sprintf("/data/%s", dataID)

	resp, err := c.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	return result, nil
}

// PostData sends one sample resource to External API 1.
func (c *APIClient) PostData(ctx context.Context, data map[string]any) (map[string]any, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	resp, err := c.Request(ctx, http.MethodPost, "/data", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	return result, nil
}
