package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// Manager defines secret retrieval behavior.
type Manager interface {
	GetSecret(ctx context.Context, secretID string) (string, error)
}

// EnvManager is a no-op secrets implementation for env-backed development.
type EnvManager struct{}

// DatabaseSecret stores database credentials loaded from Secret Manager.
type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
}

// APITokenSecret stores API token credentials loaded from Secret Manager.
type APITokenSecret struct {
	Token       string `json:"token"`
	AccessKey   string `json:"access_key,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
	BearerToken string `json:"bearer_token,omitempty"`
}

// AWSSecretsManager is a Secret Manager-backed secrets implementation.
type AWSSecretsManager struct {
	client *secretsmanager.Client
	config *config.Config
}

// NewAWSSecretsManager creates an AWS Secrets Manager client wrapper.
func NewAWSSecretsManager(cfg *config.Config) (*AWSSecretsManager, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.AWS.Region))
	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %w", err)
	}

	client := secretsmanager.NewFromConfig(awsCfg)

	return &AWSSecretsManager{
		client: client,
		config: cfg,
	}, nil
}

// NewManager creates the secrets manager configured for the current environment.
func NewManager(cfg *config.Config) (Manager, error) {
	if !cfg.UseAWSSecretsManager() {
		return EnvManager{}, nil
	}

	return NewAWSSecretsManager(cfg)
}

// GetSecret reports that env-backed mode does not resolve arbitrary secret IDs.
func (EnvManager) GetSecret(_ context.Context, secretID string) (string, error) {
	return "", fmt.Errorf("env providerではsecret idを解決できません: %s", secretID)
}

// GetSecret fetches one raw secret string by ID.
func (sm *AWSSecretsManager) GetSecret(ctx context.Context, secretID string) (string, error) {
	if !sm.config.UseAWSSecretsManager() {
		return "", fmt.Errorf("aws secrets managerは無効に設定されています")
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	result, err := sm.client.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("シークレットの取得に失敗しました: %w", err)
	}

	return *result.SecretString, nil
}

// GetDatabaseSecret fetches and parses one database secret.
func (sm *AWSSecretsManager) GetDatabaseSecret(ctx context.Context, secretID string) (*DatabaseSecret, error) {
	secretValue, err := sm.GetSecret(ctx, secretID)
	if err != nil {
		return nil, err
	}

	var dbSecret DatabaseSecret
	if err := json.Unmarshal([]byte(secretValue), &dbSecret); err != nil {
		return nil, fmt.Errorf("データベースシークレットのパースに失敗しました: %w", err)
	}

	return &dbSecret, nil
}

// GetAPIToken fetches and parses one API token secret.
func (sm *AWSSecretsManager) GetAPIToken(ctx context.Context, secretID string) (*APITokenSecret, error) {
	secretValue, err := sm.GetSecret(ctx, secretID)
	if err != nil {
		return nil, err
	}

	var tokenSecret APITokenSecret
	if err := json.Unmarshal([]byte(secretValue), &tokenSecret); err != nil {
		return nil, fmt.Errorf("APIトークンシークレットのパースに失敗しました: %w", err)
	}

	return &tokenSecret, nil
}

// FormatDSN converts the database secret into a DSN for the given dialect.
func (s *DatabaseSecret) FormatDSN(dialect string) string {
	switch dialect {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			s.Username, s.Password, s.Host, s.Port, s.DBName)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			s.Host, s.Port, s.Username, s.Password, s.DBName)
	default:
		slog.Warn("未対応のデータベースダイアレクトです", "dialect", dialect)
		return ""
	}
}
