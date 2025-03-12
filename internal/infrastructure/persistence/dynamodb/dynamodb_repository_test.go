package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetItem(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// テストデータ
	testItem := map[string]types.AttributeValue{
		"PK":   &types.AttributeValueMemberS{Value: "TEST#1"},
		"SK":   &types.AttributeValueMemberS{Value: "DETAIL#1"},
		"Name": &types.AttributeValueMemberS{Value: "Test Item"},
		"Age":  &types.AttributeValueMemberN{Value: "30"},
	}

	// モックの期待値を設定
	mockClient.On("GetItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
		return *input.TableName == "test-table" &&
			input.Key["PK"].(*types.AttributeValueMemberS).Value == "TEST#1" &&
			input.Key["SK"].(*types.AttributeValueMemberS).Value == "DETAIL#1"
	}), mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: testItem,
	}, nil)

	// テスト実行
	result, err := repo.GetItem(context.Background(), "TEST#1", "DETAIL#1")

	// アサーション
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Item", result["Name"])

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}

func TestPutItem(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// テストデータ
	testItem := map[string]interface{}{
		"PK":   "TEST#1",
		"SK":   "DETAIL#1",
		"Name": "Test Item",
		"Age":  30,
	}

	// モックの期待値を設定
	mockClient.On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
		return *input.TableName == "test-table"
	}), mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	// テスト実行
	err := repo.PutItem(context.Background(), testItem)

	// アサーション
	assert.NoError(t, err)

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// モックの期待値を設定
	mockClient.On("DeleteItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
		return *input.TableName == "test-table" &&
			input.Key["PK"].(*types.AttributeValueMemberS).Value == "TEST#1" &&
			input.Key["SK"].(*types.AttributeValueMemberS).Value == "DETAIL#1"
	}), mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)

	// テスト実行
	err := repo.DeleteItem(context.Background(), "TEST#1", "DETAIL#1")

	// アサーション
	assert.NoError(t, err)

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}

func TestQuery(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// テストデータ
	testItems := []map[string]types.AttributeValue{
		{
			"PK":   &types.AttributeValueMemberS{Value: "TEST#1"},
			"SK":   &types.AttributeValueMemberS{Value: "DETAIL#1"},
			"Name": &types.AttributeValueMemberS{Value: "Test Item 1"},
		},
		{
			"PK":   &types.AttributeValueMemberS{Value: "TEST#1"},
			"SK":   &types.AttributeValueMemberS{Value: "DETAIL#2"},
			"Name": &types.AttributeValueMemberS{Value: "Test Item 2"},
		},
	}

	// モックの期待値を設定
	mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
		return *input.TableName == "test-table"
	}), mock.Anything).Return(&dynamodb.QueryOutput{
		Items: testItems,
	}, nil)

	// テスト実行
	results, err := repo.Query(context.Background(), "TEST#1", "")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Test Item 1", results[0]["Name"])
	assert.Equal(t, "Test Item 2", results[1]["Name"])

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}

func TestScan(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// テストデータ
	testItems := []map[string]types.AttributeValue{
		{
			"PK":   &types.AttributeValueMemberS{Value: "TEST#1"},
			"SK":   &types.AttributeValueMemberS{Value: "DETAIL#1"},
			"Name": &types.AttributeValueMemberS{Value: "Test Item 1"},
		},
		{
			"PK":   &types.AttributeValueMemberS{Value: "TEST#2"},
			"SK":   &types.AttributeValueMemberS{Value: "DETAIL#1"},
			"Name": &types.AttributeValueMemberS{Value: "Test Item 2"},
		},
	}

	// モックの期待値を設定
	mockClient.On("Scan", mock.Anything, mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
		return *input.TableName == "test-table"
	}), mock.Anything).Return(&dynamodb.ScanOutput{
		Items: testItems,
	}, nil)

	// テスト実行
	results, err := repo.Scan(context.Background(), "")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Test Item 1", results[0]["Name"])
	assert.Equal(t, "Test Item 2", results[1]["Name"])

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}

func TestBatchWrite(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// テストデータ
	testItems := []map[string]interface{}{
		{
			"PK":   "TEST#1",
			"SK":   "DETAIL#1",
			"Name": "Test Item 1",
		},
		{
			"PK":   "TEST#2",
			"SK":   "DETAIL#1",
			"Name": "Test Item 2",
		},
	}

	// モックの期待値を設定
	mockClient.On("BatchWriteItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.BatchWriteItemInput) bool {
		return len(input.RequestItems["test-table"]) == 2
	}), mock.Anything).Return(&dynamodb.BatchWriteItemOutput{}, nil)

	// テスト実行
	err := repo.BatchWrite(context.Background(), testItems)

	// アサーション
	assert.NoError(t, err)

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}

func TestTransactWrite(t *testing.T) {
	// モックの作成
	mockClient := new(MockDynamoDBClient)
	repo := NewDynamoDBRepository(mockClient, "test-table")

	// テストデータ
	testOperations := []map[string]interface{}{
		{
			"OperationType": "Put",
			"PK":            "TEST#1",
			"SK":            "DETAIL#1",
			"Name":          "Test Item 1",
		},
		{
			"OperationType": "Delete",
			"PK":            "TEST#2",
			"SK":            "DETAIL#1",
		},
	}

	// モックの期待値を設定
	mockClient.On("TransactWriteItems", mock.Anything, mock.MatchedBy(func(input *dynamodb.TransactWriteItemsInput) bool {
		return len(input.TransactItems) == 2
	}), mock.Anything).Return(&dynamodb.TransactWriteItemsOutput{}, nil)

	// テスト実行
	err := repo.TransactWrite(context.Background(), testOperations)

	// アサーション
	assert.NoError(t, err)

	// モックが期待通り呼ばれたことを確認
	mockClient.AssertExpectations(t)
}
