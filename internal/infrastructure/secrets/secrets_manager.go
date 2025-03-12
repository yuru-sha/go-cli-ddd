package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// Manager はシークレット管理のインターフェースです
type Manager interface {
	// GetSecret はシークレットを取得します
	GetSecret(ctx context.Context, secretID string) (string, error)
}

// DatabaseSecret はデータベース接続情報を表します
type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
}

// APITokenSecret はAPI認証用のトークン情報を表します
type APITokenSecret struct {
	Token       string `json:"token"`
	AccessKey   string `json:"access_key,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
	BearerToken string `json:"bearer_token,omitempty"`
}

// AWSSecretsManager はAWS Secrets Managerを使用してシークレットを管理します
type AWSSecretsManager struct {
	client *secretsmanager.Client
	config *config.Config
}

// NewAWSSecretsManager は新しいAWSSecretsManagerインスタンスを作成します
func NewAWSSecretsManager(cfg *config.Config) (*AWSSecretsManager, error) {
	// AWS SDKの設定を読み込む
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %w", err)
	}

	// Secret Managerクライアントを作成
	client := secretsmanager.NewFromConfig(awsCfg)

	return &AWSSecretsManager{
		client: client,
		config: cfg,
	}, nil
}

// GetSecret は指定されたシークレットIDの値を取得します
func (sm *AWSSecretsManager) GetSecret(ctx context.Context, secretID string) (string, error) {
	// シークレットが有効でない場合はエラーを返す
	if !sm.config.AWS.Secrets.Enabled {
		return "", fmt.Errorf("Secret Managerは無効に設定されています")
	}

	// シークレットを取得
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	result, err := sm.client.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("シークレットの取得に失敗しました: %w", err)
	}

	// シークレット値を返す
	secretString := *result.SecretString
	return secretString, nil
}

// GetDatabaseSecret はデータベース接続情報を取得します
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

// GetAPIToken はAPI認証用のトークンを取得します
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

// FormatDSN はデータベース接続文字列を生成します
func (s *DatabaseSecret) FormatDSN(dialect string) string {
	switch dialect {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			s.Username, s.Password, s.Host, s.Port, s.DBName)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			s.Host, s.Port, s.Username, s.Password, s.DBName)
	default:
		log.Warn().Str("dialect", dialect).Msg("未対応のデータベースダイアレクトです")
		return ""
	}
}
