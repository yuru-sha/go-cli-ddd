package externalapi2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
)

// Client は外部API2（例：Google Ads API）のクライアントです
type Client struct {
	httpClient     *http.Client
	config         *config.Config
	secretsManager secrets.Manager
	oauthConfig    *oauth2.Config
	tokenSource    oauth2.TokenSource
}

// OAuth2Secret はOAuth2認証情報を表します
type OAuth2Secret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

// NewClient は新しいClient インスタンスを作成します
func NewClient(cfg *config.Config, httpClient *http.Client, secretsManager secrets.Manager) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// OAuth2の設定
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ExternalAPI2.ClientID,
		ClientSecret: cfg.ExternalAPI2.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		Scopes: []string{
			"https://www.googleapis.com/auth/adwords",
		},
	}

	client := &Client{
		httpClient:     httpClient,
		config:         cfg,
		secretsManager: secretsManager,
		oauthConfig:    oauthConfig,
	}

	return client
}

// GetTokenSource はOAuth2トークンソースを取得します
func (c *Client) GetTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	// すでにトークンソースが初期化されている場合はそれを返す
	if c.tokenSource != nil {
		return c.tokenSource, nil
	}

	var refreshToken string

	// Secret Managerからトークンを取得
	if c.config.AWS.Secrets.Enabled && c.config.ExternalAPI2.OAuth2SecretID != "" {
		log.Info().Msg("Secret ManagerからOAuth2認証情報を取得します")

		secretValue, err := c.secretsManager.GetSecret(ctx, c.config.ExternalAPI2.OAuth2SecretID)
		if err != nil {
			return nil, fmt.Errorf("OAuth2認証情報の取得に失敗しました: %w", err)
		}

		var oauthSecret OAuth2Secret
		if err := json.Unmarshal([]byte(secretValue), &oauthSecret); err != nil {
			return nil, fmt.Errorf("OAuth2認証情報の解析に失敗しました: %w", err)
		}

		// Secret Managerから取得した値で設定を上書き
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
		// 設定ファイルから取得
		refreshToken = c.config.ExternalAPI2.RefreshToken
	}

	if refreshToken == "" {
		return nil, fmt.Errorf("リフレッシュトークンが設定されていません")
	}

	// リフレッシュトークンからトークンを作成
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	// トークンソースを作成して保存
	c.tokenSource = c.oauthConfig.TokenSource(ctx, token)
	return c.tokenSource, nil
}

// GetAuthenticatedClient は認証済みのHTTPクライアントを取得します
func (c *Client) GetAuthenticatedClient(ctx context.Context) (*http.Client, error) {
	tokenSource, err := c.GetTokenSource(ctx)
	if err != nil {
		return nil, err
	}

	// OAuth2認証済みのHTTPクライアントを作成
	return oauth2.NewClient(ctx, tokenSource), nil
}

// Request は外部API2へのリクエストを実行します
func (c *Client) Request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.config.ExternalAPI2.BaseURL, path)

	// 認証済みのHTTPクライアントを取得
	client, err := c.GetAuthenticatedClient(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	// Google Ads APIの場合は追加のヘッダーが必要
	if c.config.ExternalAPI2.DeveloperToken != "" {
		req.Header.Set("developer-token", c.config.ExternalAPI2.DeveloperToken)
	}
	if c.config.ExternalAPI2.LoginCustomerID != "" {
		req.Header.Set("login-customer-id", c.config.ExternalAPI2.LoginCustomerID)
	}

	req.Header.Set("Content-Type", "application/json")

	// リクエストを実行
	resp, err := client.Do(req)
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

// GetCampaigns はキャンペーン情報を取得するサンプルメソッドです
func (c *Client) GetCampaigns(ctx context.Context, customerID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/customers/%s/campaigns", customerID)

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

// CreateCampaign はキャンペーンを作成するサンプルメソッドです
func (c *Client) CreateCampaign(ctx context.Context, customerID string, campaign map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/customers/%s/campaigns", customerID)

	jsonData, err := json.Marshal(campaign)
	if err != nil {
		return nil, fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	resp, err := c.Request(ctx, http.MethodPost, path, strings.NewReader(string(jsonData)))
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
