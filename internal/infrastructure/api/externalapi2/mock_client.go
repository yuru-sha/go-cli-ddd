package externalapi2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// MockClient は外部API2のモッククライアントです
type MockClient struct {
	// モックデータを保持するマップ
	mockData map[string]interface{}
	// 呼び出されたメソッドを記録
	calls []string
	// モックトークンソース
	mockTokenSource oauth2.TokenSource
}

// NewMockClient は新しいMockClientインスタンスを作成します
func NewMockClient() *MockClient {
	return &MockClient{
		mockData: make(map[string]interface{}),
		calls:    []string{},
	}
}

// SetMockData はモックデータを設定します
func (m *MockClient) SetMockData(key string, data interface{}) {
	m.mockData[key] = data
}

// GetCalls は呼び出されたメソッドのリストを返します
func (m *MockClient) GetCalls() []string {
	return m.calls
}

// recordCall はメソッド呼び出しを記録します
func (m *MockClient) recordCall(method string) {
	m.calls = append(m.calls, method)
}

// GetTokenSource はモックトークンソースを返します
func (m *MockClient) GetTokenSource(_ context.Context) (oauth2.TokenSource, error) {
	m.recordCall("GetTokenSource")

	// モックトークンソースが設定されていない場合は新しく作成
	if m.mockTokenSource == nil {
		token := &oauth2.Token{
			AccessToken:  "mock-access-token",
			RefreshToken: "mock-refresh-token",
			TokenType:    "Bearer",
		}
		m.mockTokenSource = oauth2.StaticTokenSource(token)
	}

	return m.mockTokenSource, nil
}

// GetAuthenticatedClient はモック認証済みクライアントを返します
func (m *MockClient) GetAuthenticatedClient(_ context.Context) (*http.Client, error) {
	m.recordCall("GetAuthenticatedClient")

	// 実際のHTTPクライアントは返さず、nilを返す（モックでは使用しない）
	return nil, nil
}

// GetCampaigns はモックキャンペーン情報を返します
func (m *MockClient) GetCampaigns(_ context.Context, customerID string) (map[string]interface{}, error) {
	m.recordCall(fmt.Sprintf("GetCampaigns:%s", customerID))

	// モックデータが設定されている場合はそれを返す
	if data, ok := m.mockData[fmt.Sprintf("campaigns:%s", customerID)]; ok {
		if result, ok := data.(map[string]interface{}); ok {
			return result, nil
		}
	}

	// デフォルトのモックデータを返す
	return map[string]interface{}{
		"campaigns": []map[string]interface{}{
			{
				"id":        "1234567890",
				"name":      "Mock Campaign 1",
				"status":    "ENABLED",
				"budget":    100.00,
				"customer":  customerID,
				"startDate": "2023-01-01",
				"endDate":   "2023-12-31",
			},
			{
				"id":        "0987654321",
				"name":      "Mock Campaign 2",
				"status":    "PAUSED",
				"budget":    200.00,
				"customer":  customerID,
				"startDate": "2023-02-01",
				"endDate":   "2023-11-30",
			},
		},
		"nextPageToken": "",
	}, nil
}

// CreateCampaign はモックキャンペーンを作成します
func (m *MockClient) CreateCampaign(_ context.Context, customerID string, campaign map[string]interface{}) (map[string]interface{}, error) {
	m.recordCall(fmt.Sprintf("CreateCampaign:%s", customerID))

	// 入力データをJSON文字列に変換（デバッグ用）
	jsonData, _ := json.Marshal(campaign)
	dataStr := string(jsonData)

	// データのバリデーション（例：必須フィールドのチェック）
	if _, ok := campaign["name"]; !ok {
		return nil, fmt.Errorf("必須フィールド 'name' がありません")
	}

	// 成功レスポンスを返す
	return map[string]interface{}{
		"id":           "campaign-12345",
		"name":         campaign["name"],
		"status":       "ENABLED",
		"customer":     customerID,
		"resourceName": fmt.Sprintf("customers/%s/campaigns/campaign-12345", customerID),
		"request":      dataStr,
		"createTime":   "2023-01-01T00:00:00Z",
	}, nil
}

// Request はモックリクエストを実行します（実際には何もしません）
func (m *MockClient) Request(_ context.Context, method, path string, _ io.Reader) (*http.Response, error) {
	m.recordCall(fmt.Sprintf("Request:%s:%s", method, path))

	// モックでは実際のHTTPレスポンスは返さない
	return nil, fmt.Errorf("モックでは直接Requestメソッドは使用できません。代わりにGetCampaignsやCreateCampaignなどの高レベルメソッドを使用してください")
}
