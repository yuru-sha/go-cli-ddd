package externalapi1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// APIClient は外部API1のクライアントです
type APIClient struct {
	client         *http.Client
	config         *config.Config
	secretsManager secrets.Manager
}

// TokenCache はトークンをキャッシュするための構造体です
type TokenCache struct {
	Token     string
	ExpiresAt time.Time
}

// NewAPIClient は新しいAPIClientを作成します
func NewAPIClient(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) *APIClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &APIClient{
		client:         httpClient,
		config:         cfg,
		secretsManager: secretsManager,
	}
}

// GetAuthorizationHeader は認証ヘッダーを取得します
func (c *APIClient) GetAuthorizationHeader(ctx context.Context) (string, string, error) {
	// Secret Managerが有効で、トークンのSecretIDが設定されている場合
	if c.config.AWS.Secrets.Enabled && c.config.ExternalAPI1.TokenSecretID != "" {
		log.Info().Msg("Secret ManagerからAPIトークンを取得します")

		// APIトークンを取得
		tokenStr, err := c.secretsManager.GetSecret(ctx, c.config.ExternalAPI1.TokenSecretID)
		if err != nil {
			return "", "", fmt.Errorf("APIトークンの取得に失敗しました: %w", err)
		}

		// JSONをパース
		var tokenSecret secrets.APITokenSecret
		if err := json.Unmarshal([]byte(tokenStr), &tokenSecret); err != nil {
			return "", "", fmt.Errorf("APIトークンのパースに失敗しました: %w", err)
		}

		// トークンの種類に応じて適切なヘッダーを返す
		if tokenSecret.BearerToken != "" {
			return "Authorization", fmt.Sprintf("Bearer %s", tokenSecret.BearerToken), nil
		} else if tokenSecret.Token != "" {
			return "X-API-Token", tokenSecret.Token, nil
		} else if tokenSecret.AccessKey != "" && tokenSecret.SecretKey != "" {
			// アクセスキーとシークレットキーを使用した認証（必要に応じて実装）
			return "X-API-Key", tokenSecret.AccessKey, nil
		}
	}

	// デフォルトのトークン（開発用）
	return "X-API-Token", "dev-token-12345", nil
}

// CreateAuthenticatedRequest は認証情報を含むHTTPリクエストを作成します
func (c *APIClient) CreateAuthenticatedRequest(ctx context.Context, method, url string, _ interface{}) (*http.Request, error) {
	// リクエストを作成
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	// 認証ヘッダーを取得して設定
	headerName, headerValue, err := c.GetAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if headerName != "" && headerValue != "" {
		req.Header.Set(headerName, headerValue)
	}

	// Content-Typeを設定
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Request は外部API1へのリクエストを実行します
func (c *APIClient) Request(ctx context.Context, method, path string, _ io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.config.ExternalAPI1.BaseURL, path)

	req, err := c.CreateAuthenticatedRequest(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// リクエストを実行
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの実行に失敗しました: %w", err)
	}

	// エラーレスポンスの場合
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: %s - %s", resp.Status, string(body))
	}

	return resp, nil
}

// GetData は外部API1からデータを取得するサンプルメソッドです
func (c *APIClient) GetData(ctx context.Context, dataID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/data/%s", dataID)

	resp, err := c.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	return result, nil
}

// PostData は外部API1にデータを送信するサンプルメソッドです
func (c *APIClient) PostData(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	resp, err := c.Request(ctx, http.MethodPost, "/data", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
	}

	return result, nil
}
