package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
)

// Client はDynamoDBクライアントのインターフェースです
// テスト時にモックできるようにするためのインターフェース
type Client interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

// RepositoryImpl はDynamoDBリポジトリの実装です
type RepositoryImpl struct {
	client    Client
	tableName string
}

// NewDynamoDBRepository は新しいDynamoDBリポジトリを作成します
func NewDynamoDBRepository(client Client, tableName string) repository.DynamoDBRepository {
	return &RepositoryImpl{
		client:    client,
		tableName: tableName,
	}
}

// GetItem は指定されたキーでアイテムを取得します
func (r *RepositoryImpl) GetItem(ctx context.Context, partitionKey string, sortKey string) (map[string]interface{}, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: partitionKey},
			"SK": &types.AttributeValueMemberS{Value: sortKey},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("DynamoDB GetItem error: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	item := make(map[string]interface{})
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB item: %w", err)
	}

	return item, nil
}

// PutItem は新しいアイテムを作成または更新します
func (r *RepositoryImpl) PutItem(ctx context.Context, item map[string]interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal DynamoDB item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("DynamoDB PutItem error: %w", err)
	}

	return nil
}

// DeleteItem は指定されたキーのアイテムを削除します
func (r *RepositoryImpl) DeleteItem(ctx context.Context, partitionKey string, sortKey string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: partitionKey},
			"SK": &types.AttributeValueMemberS{Value: sortKey},
		},
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("DynamoDB DeleteItem error: %w", err)
	}

	return nil
}

// Query はパーティションキーと条件に基づいてアイテムを検索します
func (r *RepositoryImpl) Query(ctx context.Context, partitionKey string, filterExpressionStr string) ([]map[string]interface{}, error) {
	// キー条件式を作成
	keyCond := expression.Key("PK").Equal(expression.Value(partitionKey))

	// フィルター式がある場合は追加
	var builder expression.Builder
	if filterExpressionStr != "" {
		// 注意: 実際のアプリケーションでは、文字列からフィルター式を構築するのではなく
		// expression.Nameとexpression.Valueを使用して安全に構築することをお勧めします
		// ここでは簡略化のため、文字列をそのまま使用しています
		builder = expression.NewBuilder().WithKeyCondition(keyCond)
		// 実際のアプリケーションでは、ここでフィルター式を適切に構築する必要があります
	} else {
		builder = expression.NewBuilder().WithKeyCondition(keyCond)
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// フィルター式がある場合は追加
	if filterExpressionStr != "" {
		input.FilterExpression = aws.String(filterExpressionStr)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("DynamoDB Query error: %w", err)
	}

	items := make([]map[string]interface{}, 0, len(result.Items))
	for _, item := range result.Items {
		m := make(map[string]interface{})
		err = attributevalue.UnmarshalMap(item, &m)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal DynamoDB item: %w", err)
		}
		items = append(items, m)
	}

	return items, nil
}

// Scan はテーブル全体をスキャンして条件に一致するアイテムを検索します
func (r *RepositoryImpl) Scan(ctx context.Context, filterExpressionStr string) ([]map[string]interface{}, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	// フィルター式がある場合は追加
	if filterExpressionStr != "" {
		input.FilterExpression = aws.String(filterExpressionStr)
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("DynamoDB Scan error: %w", err)
	}

	items := make([]map[string]interface{}, 0, len(result.Items))
	for _, item := range result.Items {
		m := make(map[string]interface{})
		err = attributevalue.UnmarshalMap(item, &m)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal DynamoDB item: %w", err)
		}
		items = append(items, m)
	}

	return items, nil
}

// BatchWrite は複数のアイテムを一括で書き込みます
func (r *RepositoryImpl) BatchWrite(ctx context.Context, items []map[string]interface{}) error {
	if len(items) == 0 {
		return nil
	}

	// DynamoDBのBatchWriteItemは一度に25項目までしか処理できないため、
	// 25項目ごとにバッチを分割する必要があります
	const maxBatchSize = 25
	for i := 0; i < len(items); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		writeRequests := make([]types.WriteRequest, len(batch))

		for j, item := range batch {
			av, err := attributevalue.MarshalMap(item)
			if err != nil {
				return fmt.Errorf("failed to marshal DynamoDB item: %w", err)
			}

			writeRequests[j] = types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: av,
				},
			}
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests,
			},
		}

		_, err := r.client.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("DynamoDB BatchWriteItem error: %w", err)
		}
	}

	return nil
}

// TransactWrite はトランザクション内で複数の書き込み操作を実行します
func (r *RepositoryImpl) TransactWrite(ctx context.Context, operations []map[string]interface{}) error {
	if len(operations) == 0 {
		return nil
	}

	// DynamoDBのTransactWriteItemsは一度に100項目までしか処理できないため、
	// 100項目以上の場合はエラーを返します
	const maxTransactSize = 100
	if len(operations) > maxTransactSize {
		return errors.New("too many operations for a single transaction (max 100)")
	}

	transactItems := make([]types.TransactWriteItem, len(operations))

	for i, op := range operations {
		// 操作タイプを取得（Put, Update, Delete, ConditionCheck）
		opType, ok := op["OperationType"].(string)
		if !ok {
			return fmt.Errorf("operation type not specified for item %d", i)
		}

		delete(op, "OperationType")

		switch opType {
		case "Put":
			av, err := attributevalue.MarshalMap(op)
			if err != nil {
				return fmt.Errorf("failed to marshal DynamoDB item: %w", err)
			}

			transactItems[i] = types.TransactWriteItem{
				Put: &types.Put{
					TableName: aws.String(r.tableName),
					Item:      av,
				},
			}

		case "Delete":
			pk, pkOk := op["PK"].(string)
			sk, skOk := op["SK"].(string)
			if !pkOk || !skOk {
				return fmt.Errorf("PK or SK not specified for Delete operation %d", i)
			}

			transactItems[i] = types.TransactWriteItem{
				Delete: &types.Delete{
					TableName: aws.String(r.tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: pk},
						"SK": &types.AttributeValueMemberS{Value: sk},
					},
				},
			}

		default:
			return fmt.Errorf("unsupported operation type: %s", opType)
		}
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}

	_, err := r.client.TransactWriteItems(ctx, input)
	if err != nil {
		return fmt.Errorf("DynamoDB TransactWriteItems error: %w", err)
	}

	return nil
}
