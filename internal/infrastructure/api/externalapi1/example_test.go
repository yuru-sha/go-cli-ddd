package externalapi1

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

	// GetTokenのテスト
	token, err := mockClient.GetToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "mock-api1-token", token)

	// GetDataのテスト
	data, err := mockClient.GetData(ctx, "test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", data["id"])
	assert.Equal(t, "Mock Data test-id", data["name"])

	// カスタムモックデータを設定
	customData := map[string]interface{}{
		"id":      "custom-id",
		"name":    "Custom Data",
		"custom":  true,
		"version": 1,
	}
	mockClient.SetMockData("data:custom-id", customData)

	// カスタムモックデータを取得
	data, err = mockClient.GetData(ctx, "custom-id")
	assert.NoError(t, err)
	assert.Equal(t, "custom-id", data["id"])
	assert.Equal(t, "Custom Data", data["name"])
	assert.Equal(t, true, data["custom"])
	assert.Equal(t, int(1), data["version"])

	// PostDataのテスト
	postData := map[string]interface{}{
		"name":        "Test Post",
		"description": "This is a test post",
		"tags":        []string{"test", "mock"},
	}

	result, err := mockClient.PostData(ctx, postData)
	assert.NoError(t, err)
	assert.Equal(t, "mock-id-12345", result["id"])
	assert.Equal(t, postData["name"], result["name"])
	assert.Equal(t, "created", result["status"])

	// 呼び出し履歴の確認
	calls := mockClient.GetCalls()
	assert.Contains(t, calls, "GetToken")
	assert.Contains(t, calls, "GetData:test-id")
	assert.Contains(t, calls, "GetData:custom-id")
	assert.Contains(t, calls, "PostData")
}

func TestClientConfig(t *testing.T) {
	// このテストはモックを使用しないため、実際のAPIを呼び出さないようにスキップ
	t.Skip("This test requires actual API credentials")

	// 設定を作成
	cfg := &config.Config{
		ExternalAPI1: config.ExternalAPI1Config{
			BaseURL:       "https://api.example.com/v1",
			TokenSecretID: "api1/token",
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
	client := NewAPIClient(cfg, nil, nil)

	// クライアントの設定が正しく行われていることを確認
	assert.Equal(t, "https://api.example.com/v1", client.config.ExternalAPI1.BaseURL)
}
