package externalapi1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// MockClient は外部API1のモッククライアントです
type MockClient struct {
	// モックデータを保持するマップ
	mockData map[string]any
	// 呼び出されたメソッドを記録
	calls []string
}

// NewMockClient は新しいMockClientインスタンスを作成します
func NewMockClient() *MockClient {
	return &MockClient{
		mockData: make(map[string]any),
		calls:    []string{},
	}
}

// SetMockData はモックデータを設定します
func (m *MockClient) SetMockData(key string, data any) {
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

// GetToken はモックトークンを返します
func (m *MockClient) GetToken(_ context.Context) (string, error) {
	m.recordCall("GetToken")
	return "mock-api1-token", nil
}

// GetData はモックデータを返します
func (m *MockClient) GetData(_ context.Context, dataID string) (map[string]any, error) {
	m.recordCall(fmt.Sprintf("GetData:%s", dataID))

	// モックデータが設定されている場合はそれを返す
	if data, ok := m.mockData[fmt.Sprintf("data:%s", dataID)]; ok {
		if result, ok := data.(map[string]any); ok {
			return result, nil
		}
	}

	// デフォルトのモックデータを返す
	return map[string]any{
		"id":         dataID,
		"name":       fmt.Sprintf("Mock Data %s", dataID),
		"type":       "mock",
		"created_at": "2023-01-01T00:00:00Z",
	}, nil
}

// PostData はモックデータを受け取り、処理したように見せかけます
func (m *MockClient) PostData(_ context.Context, data map[string]any) (map[string]any, error) {
	m.recordCall("PostData")

	// 入力データをJSON文字列に変換（デバッグ用）
	jsonData, _ := json.Marshal(data)
	dataStr := string(jsonData)

	// データのバリデーション（例：必須フィールドのチェック）
	if _, ok := data["name"]; !ok {
		return nil, fmt.Errorf("必須フィールド 'name' がありません")
	}

	// 成功レスポンスを返す
	return map[string]any{
		"id":        "mock-id-12345",
		"name":      data["name"],
		"status":    "created",
		"request":   dataStr,
		"timestamp": "2023-01-01T00:00:00Z",
	}, nil
}

// Request はモックリクエストを実行します（実際には何もしません）
func (m *MockClient) Request(ctx context.Context, method, path string, body io.Reader) (any, error) {
	m.recordCall(fmt.Sprintf("Request:%s:%s", method, path))

	// リクエストボディを読み込む（あれば）
	var bodyStr string
	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("リクエストボディの読み込みに失敗しました: %w", err)
		}
		bodyStr = string(bodyBytes)
	}

	// パスに基づいてモックレスポンスを返す
	if strings.Contains(path, "/data/") && method == "GET" {
		// データID抽出
		parts := strings.Split(path, "/")
		dataID := parts[len(parts)-1]

		return m.GetData(ctx, dataID)
	} else if path == "/data" && method == "POST" {
		// POSTリクエストの場合
		var data map[string]any
		if err := json.Unmarshal([]byte(bodyStr), &data); err != nil {
			return nil, fmt.Errorf("リクエストボディのパースに失敗しました: %w", err)
		}

		return m.PostData(ctx, data)
	}

	// 未対応のパスの場合
	return map[string]any{
		"status": "mock_response",
		"path":   path,
		"method": method,
		"body":   bodyStr,
	}, nil
}
