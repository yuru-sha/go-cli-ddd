package externalapi2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// Client は外部API2（例：Google Ads API）のクライアントです。
type Client struct {
	httpClient     *http.Client
	config         *config.Config
	secretsManager secrets.Manager
	oauthConfig    *oauth2.Config
	tokenSource    oauth2.TokenSource
}

// OAuth2Secret stores OAuth2 credentials loaded from Secret Manager.
type OAuth2Secret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

// NewClient creates an External API 2 client.
func NewClient(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ExternalAPI2.ClientID,
		ClientSecret: cfg.ExternalAPI2.ClientSecret,
		Endpoint: oauth2.Endpoint{ //nolint:gosec // Fixed Google OAuth endpoints are public URLs, not credentials.
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		Scopes: []string{"https://www.googleapis.com/auth/adwords"},
	}

	return &Client{
		httpClient:     httpClient,
		config:         cfg,
		secretsManager: secretsManager,
		oauthConfig:    oauthConfig,
	}
}

// GetTokenSource resolves the OAuth2 token source for External API 2.
func (c *Client) GetTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	if c.tokenSource != nil {
		return c.tokenSource, nil
	}

	var refreshToken string
	if c.config.UseAWSSecretsManager() && c.config.ExternalAPI2.OAuth2SecretID != "" {
		slog.Info("Secret ManagerからOAuth2認証情報を取得します")

		secretValue, err := c.secretsManager.GetSecret(ctx, c.config.ExternalAPI2.OAuth2SecretID)
		if err != nil {
			return nil, fmt.Errorf("OAuth2認証情報の取得に失敗しました: %w", err)
		}

		var oauthSecret OAuth2Secret
		if err := json.Unmarshal([]byte(secretValue), &oauthSecret); err != nil {
			return nil, fmt.Errorf("OAuth2認証情報の解析に失敗しました: %w", err)
		}

		if oauthSecret.ClientID != "" {
			c.oauthConfig.ClientID = oauthSecret.ClientID
		}
		if oauthSecret.ClientSecret != "" {
			c.oauthConfig.ClientSecret = oauthSecret.ClientSecret
		}
		if oauthSecret.RefreshToken != "" {
			refreshToken = oauthSecret.RefreshToken
		}
	} else {
		refreshToken = c.config.ExternalAPI2.RefreshToken
	}

	if refreshToken == "" {
		return nil, fmt.Errorf("リフレッシュトークンが設定されていません")
	}

	token := &oauth2.Token{RefreshToken: refreshToken}
	c.tokenSource = c.oauthConfig.TokenSource(ctx, token)
	return c.tokenSource, nil
}

// GetAuthenticatedClient returns an authenticated HTTP client for External API 2.
func (c *Client) GetAuthenticatedClient(ctx context.Context) (*http.Client, error) {
	tokenSource, err := c.GetTokenSource(ctx)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(ctx, tokenSource), nil
}

// Request executes a request against External API 2.
func (c *Client) Request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.config.ExternalAPI2.BaseURL, path)

	client, err := c.GetAuthenticatedClient(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	if c.config.ExternalAPI2.DeveloperToken != "" {
		req.Header.Set("developer-token", c.config.ExternalAPI2.DeveloperToken)
	}
	if c.config.ExternalAPI2.LoginCustomerID != "" {
		req.Header.Set("login-customer-id", c.config.ExternalAPI2.LoginCustomerID)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

// GetCampaigns fetches campaigns for one customer from External API 2.
func (c *Client) GetCampaigns(ctx context.Context, customerID string) (map[string]any, error) {
	path := fmt.Sprintf("/customers/%s/campaigns", customerID)

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

// CreateCampaign creates a campaign for one customer in External API 2.
func (c *Client) CreateCampaign(ctx context.Context, customerID string, campaign map[string]any) (map[string]any, error) {
	path := fmt.Sprintf("/customers/%s/campaigns", customerID)

	jsonData, err := json.Marshal(campaign)
	if err != nil {
		return nil, fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	resp, err := c.Request(ctx, http.MethodPost, path, strings.NewReader(string(jsonData)))
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
