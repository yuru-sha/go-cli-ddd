package config

import (
	"fmt"
	"strings"

	"dario.cat/mergo"
	"github.com/spf13/viper"
)

// Config はアプリケーション設定を管理します
type Config struct {
	App          AppConfig          `mapstructure:"app"`
	Database     DatabaseConfig     `mapstructure:"database"`
	HTTP         HTTPConfig         `mapstructure:"http"`
	AWS          AWSConfig          `mapstructure:"aws"`
	Notification NotificationConfig `mapstructure:"notification"`
	ExternalAPI1 ExternalAPI1Config `mapstructure:"external_api1"`
	ExternalAPI2 ExternalAPI2Config `mapstructure:"external_api2"`
}

// AppConfig はアプリケーション全般の設定です
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

// DatabaseConfig はデータベース接続の設定です
type DatabaseConfig struct {
	Dialect     string       `mapstructure:"dialect"`
	DSN         string       `mapstructure:"dsn"`
	LogLevel    string       `mapstructure:"log_level"`
	AutoMigrate bool         `mapstructure:"auto_migrate"`
	SecretID    string       `mapstructure:"secret_id"`
	Aurora      AuroraConfig `mapstructure:"aurora"`
}

// AuroraConfig はAWS Aurora接続の設定です
type AuroraConfig struct {
	Enabled bool               `mapstructure:"enabled"`
	Writer  AuroraNodeConfig   `mapstructure:"writer"`
	Reader  AuroraReaderConfig `mapstructure:"reader"`
}

// AuroraNodeConfig はAurora単一ノードの設定です
type AuroraNodeConfig struct {
	SecretID string `mapstructure:"secret_id"`
}

// AuroraReaderConfig はAuroraリードレプリカの設定です
type AuroraReaderConfig struct {
	SecretID      string `mapstructure:"secret_id"`
	LoadBalancing string `mapstructure:"load_balancing"` // "random" または "round-robin"
}

// HTTPConfig はHTTPクライアントの設定です
type HTTPConfig struct {
	Timeout    int             `mapstructure:"timeout"`
	MaxRetries int             `mapstructure:"max_retries"`
	RateLimit  RateLimitConfig `mapstructure:"rate_limit"`
}

// RateLimitConfig はレート制限の設定です
type RateLimitConfig struct {
	QPS   float64 `mapstructure:"qps"`
	Burst int     `mapstructure:"burst"`
}

// NotificationConfig は通知関連の設定です
type NotificationConfig struct {
	Slack SlackConfig `mapstructure:"slack"`
}

// SlackConfig はSlack通知の設定です
type SlackConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	WebhookURL      string `mapstructure:"webhook_url"`
	WebhookSecretID string `mapstructure:"webhook_secret_id"`
	Channel         string `mapstructure:"channel"`
	Username        string `mapstructure:"username"`
	IconEmoji       string `mapstructure:"icon_emoji"`
	SuccessEmoji    string `mapstructure:"success_emoji"`
	FailureEmoji    string `mapstructure:"failure_emoji"`
}

// AWSConfig はAWS関連の設定です
type AWSConfig struct {
	Region  string        `mapstructure:"region"`
	Secrets SecretsConfig `mapstructure:"secrets"`
}

// SecretsConfig はAWS Secret Manager関連の設定です
type SecretsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// ExternalAPI1Config は外部API1（例：DOMO API）の設定です
type ExternalAPI1Config struct {
	BaseURL       string `mapstructure:"base_url"`
	TokenSecretID string `mapstructure:"token_secret_id"`
}

// ExternalAPI2Config は外部API2（例：Google Ads API）の設定です
type ExternalAPI2Config struct {
	BaseURL         string `mapstructure:"base_url"`
	ClientID        string `mapstructure:"client_id"`
	ClientSecret    string `mapstructure:"client_secret"`
	RefreshToken    string `mapstructure:"refresh_token"`
	OAuth2SecretID  string `mapstructure:"oauth2_secret_id"`
	DeveloperToken  string `mapstructure:"developer_token"`
	LoginCustomerID string `mapstructure:"login_customer_id"`
}

// Options は設定読み込みのオプションを表します
type Options struct {
	ConfigPath string
	Env        string
}

// NewConfigOptions は設定読み込みのオプションを作成します
func NewConfigOptions(configPath, env string) *Options {
	return &Options{
		ConfigPath: configPath,
		Env:        env,
	}
}

// LoadConfig は設定ファイルから設定を読み込みます
func LoadConfig(opts *Options) (*Config, error) {
	viper.SetConfigFile(opts.ConfigPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	// 環境が指定されていない場合はlocalを使用
	env := opts.Env
	if env == "" {
		env = "local"
	}

	// ベース環境は常にlocal
	baseEnv := "local"

	// ベース環境と指定環境が同じ場合は、そのまま設定を返す
	if strings.EqualFold(baseEnv, env) {
		var config Config
		if err := viper.UnmarshalKey(strings.ToLower(env), &config); err != nil {
			return nil, fmt.Errorf("設定の解析に失敗しました: %w", err)
		}
		return &config, nil
	}

	// ベース環境の設定を取得
	var baseConfig Config
	if err := viper.UnmarshalKey(strings.ToLower(baseEnv), &baseConfig); err != nil {
		return nil, fmt.Errorf("ベース設定の解析に失敗しました: %w", err)
	}

	// 指定された環境の設定を取得
	var envConfig Config
	if err := viper.UnmarshalKey(strings.ToLower(env), &envConfig); err != nil {
		return nil, fmt.Errorf("環境設定の解析に失敗しました: %w", err)
	}

	// ベース設定に環境設定をマージ（環境設定が優先）
	if err := mergo.Merge(&baseConfig, envConfig, mergo.WithOverride); err != nil {
		return nil, fmt.Errorf("設定のマージに失敗しました: %w", err)
	}

	return &baseConfig, nil
}
