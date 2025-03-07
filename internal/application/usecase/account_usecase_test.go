package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/entity"
)

// AccountRepositoryMock はAccountRepositoryのモック実装です
type AccountRepositoryMock struct {
	mock.Mock
}

func (m *AccountRepositoryMock) FindAll(ctx context.Context) ([]entity.Account, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Account), args.Error(1)
}

func (m *AccountRepositoryMock) FindByID(ctx context.Context, id uint) (*entity.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountRepositoryMock) Create(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *AccountRepositoryMock) Update(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *AccountRepositoryMock) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AccountRepositoryMock) SaveAll(ctx context.Context, accounts []entity.Account) error {
	args := m.Called(ctx, accounts)
	return args.Error(0)
}

// AccountAPIRepositoryMock はAccountAPIRepositoryのモック実装です
type AccountAPIRepositoryMock struct {
	mock.Mock
}

func (m *AccountAPIRepositoryMock) FetchAccounts(ctx context.Context) ([]entity.Account, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Account), args.Error(1)
}

// TestSyncAccounts はSyncAccountsメソッドのテストです
func TestSyncAccounts(t *testing.T) {
	// モックの作成
	accountRepo := new(AccountRepositoryMock)
	accountAPIRepo := new(AccountAPIRepositoryMock)

	// テスト用のアカウントデータ
	now := time.Now()
	mockAccounts := []entity.Account{
		{
			ID:        1,
			Name:      "テスト株式会社",
			Status:    "active",
			APIKey:    "test_api_key_1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "サンプル会社",
			Status:    "inactive",
			APIKey:    "test_api_key_2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// モックの振る舞いを設定
	accountAPIRepo.On("FetchAccounts", mock.Anything).Return(mockAccounts, nil)
	accountRepo.On("SaveAll", mock.Anything, mockAccounts).Return(nil)

	// テスト対象のユースケースを作成
	useCase := NewAccountUseCase(accountRepo, accountAPIRepo)

	// テスト実行
	err := useCase.SyncAccounts(context.Background())

	// アサーション
	assert.NoError(t, err)
	accountAPIRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}
