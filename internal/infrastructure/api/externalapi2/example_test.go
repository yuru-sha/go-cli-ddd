package externalapi2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

func TestMockClient(t *testing.T) {
	// モッククライアントを作成
	mockClient := NewMockClient()

	// コンテキストを作成
	ctx := context.Background()

	// GetTokenSourceのテスト
	tokenSource, err := mockClient.GetTokenSource(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tokenSource)

	// トークンの取得
	token, err := tokenSource.Token()
	assert.NoError(t, err)
	assert.Equal(t, "mock-access-token", token.AccessToken)
	assert.Equal(t, "mock-refresh-token", token.RefreshToken)
	assert.Equal(t, "Bearer", token.TokenType)

	// GetCampaignsのテスト
	customerID := "1234567890"
	campaigns, err := mockClient.GetCampaigns(ctx, customerID)
	assert.NoError(t, err)
	assert.NotNil(t, campaigns["campaigns"])

	// カスタムモックデータを設定
	customData := map[string]interface{}{
		"campaigns": []map[string]interface{}{
			{
				"id":       "custom-campaign-1",
				"name":     "Custom Campaign 1",
				"status":   "ENABLED",
				"customer": customerID,
			},
		},
		"nextPageToken": "custom-token",
	}
	mockClient.SetMockData("campaigns:"+customerID, customData)

	// カスタムモックデータを取得
	campaigns, err = mockClient.GetCampaigns(ctx, customerID)
	assert.NoError(t, err)
	assert.Equal(t, "custom-token", campaigns["nextPageToken"])
	campaignsList := campaigns["campaigns"].([]map[string]interface{})
	assert.Equal(t, 1, len(campaignsList))
	assert.Equal(t, "custom-campaign-1", campaignsList[0]["id"])
	assert.Equal(t, "Custom Campaign 1", campaignsList[0]["name"])

	// CreateCampaignのテスト
	newCampaign := map[string]interface{}{
		"name":   "New Test Campaign",
		"budget": 500.00,
		"status": "PAUSED",
	}

	result, err := mockClient.CreateCampaign(ctx, customerID, newCampaign)
	assert.NoError(t, err)
	assert.Equal(t, "campaign-12345", result["id"])
	assert.Equal(t, newCampaign["name"], result["name"])
	assert.Equal(t, "ENABLED", result["status"])
	assert.Equal(t, customerID, result["customer"])

	// 呼び出し履歴の確認
	calls := mockClient.GetCalls()
	assert.Contains(t, calls, "GetTokenSource")
	assert.Contains(t, calls, "GetCampaigns:"+customerID)
	assert.Contains(t, calls, "CreateCampaign:"+customerID)
}

func TestClientConfig(t *testing.T) {
	// このテストはモックを使用しないため、実際のAPIを呼び出さないようにスキップ
	t.Skip("This test requires actual API credentials")

	// 設定を作成
	cfg := &config.Config{
		ExternalAPI2: config.ExternalAPI2Config{
			BaseURL:         "https://googleads.googleapis.com/v14",
			ClientID:        "test-client-id",
			ClientSecret:    "test-client-secret",
			RefreshToken:    "test-refresh-token",
			OAuth2SecretID:  "googleads/oauth2",
			DeveloperToken:  "test-developer-token",
			LoginCustomerID: "1234567890",
		},
		AWS: config.AWSConfig{
			Region: "ap-northeast-1",
			Secrets: config.SecretsConfig{
				Enabled: true,
			},
		},
	}

	// HTTPクライアントとシークレットマネージャーは実際のテストでは適切なものを使用する
	// ここではnilを渡して、実際のAPIを呼び出さないようにする
	client := NewClient(cfg, nil, nil)

	// クライアントの設定が正しく行われていることを確認
	assert.Equal(t, "https://googleads.googleapis.com/v14", client.config.ExternalAPI2.BaseURL)
	assert.Equal(t, "test-client-id", client.oauthConfig.ClientID)
	assert.Equal(t, "test-client-secret", client.oauthConfig.ClientSecret)
}
