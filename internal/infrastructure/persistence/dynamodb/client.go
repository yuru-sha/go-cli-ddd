package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoDBClient はAWS SDK for Go v2を使用してDynamoDBクライアントを作成します
func NewDynamoDBClient(ctx context.Context, region string) (Client, error) {
	// AWS SDKの設定をロード
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("AWS設定のロードに失敗しました: %w", err)
	}

	// DynamoDBクライアントを作成
	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}

// NewLocalDynamoDBClient はローカル開発用のDynamoDBクライアントを作成します
func NewLocalDynamoDBClient(ctx context.Context, endpoint string) (Client, error) {
	// AWS SDKの設定をロード（リージョンはダミー値でOK）
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定のロードに失敗しました: %w", err)
	}

	// カスタムエンドポイントを使用するオプションを設定
	options := dynamodb.Options{
		EndpointResolver: dynamodb.EndpointResolverFromURL(endpoint),
	}

	// DynamoDBクライアントを作成
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		*o = options
	})
	return client, nil
}
