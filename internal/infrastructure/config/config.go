// Package config provides environment-aware configuration loading.
package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config はアプリケーション設定を管理します。
type Config struct {
	App          AppConfig          `mapstructure:"app"`
	Database     DatabaseConfig     `mapstructure:"database"`
	HTTP         HTTPConfig         `mapstructure:"http"`
	AWS          AWSConfig          `mapstructure:"aws"`
	Notification NotificationConfig `mapstructure:"notification"`
	ExternalAPI1 ExternalAPI1Config `mapstructure:"external_api1"`
	ExternalAPI2 ExternalAPI2Config `mapstructure:"external_api2"`
}

// AppConfig はアプリケーション全般の設定です。
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

// DatabaseConfig はデータベース接続の設定です。
type DatabaseConfig struct {
	Dialect     string       `mapstructure:"dialect"`
	DSN         string       `mapstructure:"dsn"`
	LogLevel    string       `mapstructure:"log_level"`
	AutoMigrate bool         `mapstructure:"auto_migrate"`
	SecretID    string       `mapstructure:"secret_id"`
	Aurora      AuroraConfig `mapstructure:"aurora"`
}

// AuroraConfig は Aurora 接続の設定です。
type AuroraConfig struct {
	Enabled bool               `mapstructure:"enabled"`
	Writer  AuroraNodeConfig   `mapstructure:"writer"`
	Reader  AuroraReaderConfig `mapstructure:"reader"`
}

// AuroraNodeConfig は Aurora ノードの設定です。
type AuroraNodeConfig struct {
	SecretID string `mapstructure:"secret_id"`
}

// AuroraReaderConfig は Aurora リーダーの設定です。
type AuroraReaderConfig struct {
	SecretID      string `mapstructure:"secret_id"`
	LoadBalancing string `mapstructure:"load_balancing"`
}

// HTTPConfig は HTTP クライアントの設定です。
type HTTPConfig struct {
	Timeout    int             `mapstructure:"timeout"`
	MaxRetries int             `mapstructure:"max_retries"`
	RateLimit  RateLimitConfig `mapstructure:"rate_limit"`
}

// RateLimitConfig はレート制限の設定です。
type RateLimitConfig struct {
	QPS   float64 `mapstructure:"qps"`
	Burst int     `mapstructure:"burst"`
}

// NotificationConfig は通知関連の設定です。
type NotificationConfig struct {
	Slack SlackConfig `mapstructure:"slack"`
}

// SlackConfig は Slack 通知の設定です。
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

// AWSConfig は AWS 関連の設定です。
type AWSConfig struct {
	Region  string        `mapstructure:"region"`
	Secrets SecretsConfig `mapstructure:"secrets"`
}

// SecretsConfig は secrets の解決方式です。
type SecretsConfig struct {
	Provider string `mapstructure:"provider"`
}

// ExternalAPI1Config は外部 API 1 の設定です。
type ExternalAPI1Config struct {
	BaseURL       string `mapstructure:"base_url"`
	Token         string `mapstructure:"token"`
	TokenSecretID string `mapstructure:"token_secret_id"`
}

// ExternalAPI2Config は外部 API 2 の設定です。
type ExternalAPI2Config struct {
	BaseURL         string `mapstructure:"base_url"`
	ClientID        string `mapstructure:"client_id"`
	ClientSecret    string `mapstructure:"client_secret"`
	RefreshToken    string `mapstructure:"refresh_token"`
	OAuth2SecretID  string `mapstructure:"oauth2_secret_id"`
	DeveloperToken  string `mapstructure:"developer_token"`
	LoginCustomerID string `mapstructure:"login_customer_id"`
}

// Options は設定読み込み時のオプションです。
type Options struct {
	ConfigPath string
	Env        string
}

// NewConfigOptions は設定読み込みオプションを作成します。
func NewConfigOptions(configPath, env string) *Options {
	return &Options{ConfigPath: configPath, Env: env}
}

// LoadConfig は設定ファイルから設定を読み込みます。
// prd を基底値にし、local では .env を最終上書きとして扱います。
func LoadConfig(opts *Options) (*Config, error) {
	bootstrapEnv, err := readDotEnvFile(".env")
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(opts.ConfigPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	env := strings.ToLower(strings.TrimSpace(opts.Env))
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	}
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(bootstrapEnv["ENV"]))
	}
	if env == "" {
		env = "prd"
	}
	switch env {
	case "local", "dev", "prd":
	default:
		return nil, fmt.Errorf("未対応の環境です: %s", env)
	}

	baseEnv := "prd"

	mergedMap := v.GetStringMap(baseEnv)

	if env != baseEnv {
		mergedMap = mergeMaps(mergedMap, v.GetStringMap(env))
	}

	var baseConfig Config
	if err := decodeConfigMap(mergedMap, &baseConfig); err != nil {
		return nil, fmt.Errorf("設定の解析に失敗しました: %w", err)
	}

	if env == "local" {
		if err := loadDotEnvMap(bootstrapEnv); err != nil {
			return nil, err
		}
	}
	applyEnvOverrides(&baseConfig)
	baseConfig.normalizeSecretsProvider()

	return &baseConfig, nil
}

// UseAWSSecretsManager reports whether secrets should be resolved from AWS Secrets Manager.
func (cfg *Config) UseAWSSecretsManager() bool {
	return cfg.AWS.Secrets.Provider == "aws"
}

func mergeMaps(base, override map[string]any) map[string]any {
	merged := make(map[string]any, len(base))
	maps.Copy(merged, base)

	for key, value := range override {
		overrideMap, overrideIsMap := value.(map[string]any)
		baseMap, baseIsMap := merged[key].(map[string]any)
		if overrideIsMap && baseIsMap {
			merged[key] = mergeMaps(baseMap, overrideMap)
			continue
		}
		merged[key] = value
	}

	return merged
}

func decodeConfigMap(data map[string]any, dest *Config) error {
	normalised, err := normalizeMap(data)
	if err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "mapstructure",
		Result:  dest,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(normalised)
}

func normalizeMap(data map[string]any) (map[string]any, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var normalized map[string]any
	if err := json.Unmarshal(bytes, &normalized); err != nil {
		return nil, err
	}

	return normalized, nil
}

func readDotEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // Local development intentionally reads the fixed .env file.
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf(".env の読み込みに失敗しました: %w", err)
	}

	values := make(map[string]string)
	for line := range strings.SplitSeq(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}

		values[key] = value
	}

	return values, nil
}

func loadDotEnvMap(values map[string]string) error {
	for key, value := range values {
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf(".env の環境変数設定に失敗しました: %w", err)
		}
	}

	return nil
}

func applyEnvOverrides(cfg *Config) {
	overrideString("APP_NAME", &cfg.App.Name)
	overrideBool("APP_DEBUG", &cfg.App.Debug)
	overrideString("APP_LOG_LEVEL", &cfg.App.LogLevel)

	overrideString("DATABASE_DIALECT", &cfg.Database.Dialect)
	overrideString("DATABASE_DSN", &cfg.Database.DSN)
	overrideString("DATABASE_LOG_LEVEL", &cfg.Database.LogLevel)
	overrideBool("DATABASE_AUTO_MIGRATE", &cfg.Database.AutoMigrate)
	overrideString("DATABASE_SECRET_ID", &cfg.Database.SecretID)
	overrideBool("DATABASE_AURORA_ENABLED", &cfg.Database.Aurora.Enabled)
	overrideString("DATABASE_AURORA_WRITER_SECRET_ID", &cfg.Database.Aurora.Writer.SecretID)
	overrideString("DATABASE_AURORA_READER_SECRET_ID", &cfg.Database.Aurora.Reader.SecretID)
	overrideString("DATABASE_AURORA_READER_LOAD_BALANCING", &cfg.Database.Aurora.Reader.LoadBalancing)

	overrideInt("HTTP_TIMEOUT", &cfg.HTTP.Timeout)
	overrideInt("HTTP_MAX_RETRIES", &cfg.HTTP.MaxRetries)
	overrideFloat("HTTP_RATE_LIMIT_QPS", &cfg.HTTP.RateLimit.QPS)
	overrideInt("HTTP_RATE_LIMIT_BURST", &cfg.HTTP.RateLimit.Burst)

	overrideString("AWS_REGION", &cfg.AWS.Region)
	overrideStringAlias([]string{"SECRETS_PROVIDER", "SECRET_PROVIDER"}, &cfg.AWS.Secrets.Provider)

	overrideString("EXTERNAL_API1_BASE_URL", &cfg.ExternalAPI1.BaseURL)
	overrideString("EXTERNAL_API1_TOKEN", &cfg.ExternalAPI1.Token)
	overrideString("EXTERNAL_API1_TOKEN_SECRET_ID", &cfg.ExternalAPI1.TokenSecretID)

	overrideString("EXTERNAL_API2_BASE_URL", &cfg.ExternalAPI2.BaseURL)
	overrideString("EXTERNAL_API2_CLIENT_ID", &cfg.ExternalAPI2.ClientID)
	overrideString("EXTERNAL_API2_CLIENT_SECRET", &cfg.ExternalAPI2.ClientSecret)
	overrideString("EXTERNAL_API2_REFRESH_TOKEN", &cfg.ExternalAPI2.RefreshToken)
	overrideString("EXTERNAL_API2_OAUTH2_SECRET_ID", &cfg.ExternalAPI2.OAuth2SecretID)
	overrideString("EXTERNAL_API2_DEVELOPER_TOKEN", &cfg.ExternalAPI2.DeveloperToken)
	overrideString("EXTERNAL_API2_LOGIN_CUSTOMER_ID", &cfg.ExternalAPI2.LoginCustomerID)

	overrideBool("SLACK_ENABLED", &cfg.Notification.Slack.Enabled)
	overrideString("SLACK_WEBHOOK_URL", &cfg.Notification.Slack.WebhookURL)
	overrideString("SLACK_WEBHOOK_SECRET_ID", &cfg.Notification.Slack.WebhookSecretID)
	overrideString("SLACK_CHANNEL", &cfg.Notification.Slack.Channel)
	overrideString("SLACK_USERNAME", &cfg.Notification.Slack.Username)
	overrideString("SLACK_ICON_EMOJI", &cfg.Notification.Slack.IconEmoji)
	overrideString("SLACK_SUCCESS_EMOJI", &cfg.Notification.Slack.SuccessEmoji)
	overrideString("SLACK_FAILURE_EMOJI", &cfg.Notification.Slack.FailureEmoji)

	if value, ok := os.LookupEnv("AWS_SECRETS_ENABLED"); ok && value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			if parsed {
				cfg.AWS.Secrets.Provider = "aws"
			} else {
				cfg.AWS.Secrets.Provider = "env"
			}
		}
	}
}

func overrideString(key string, dest *string) {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		*dest = value
	}
}

func overrideStringAlias(keys []string, dest *string) {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok && value != "" {
			*dest = value
			return
		}
	}
}

func overrideBool(key string, dest *bool) {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			*dest = parsed
		}
	}
}

func overrideInt(key string, dest *int) {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			*dest = parsed
		}
	}
}

func overrideFloat(key string, dest *float64) {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			*dest = parsed
		}
	}
}

func (cfg *Config) normalizeSecretsProvider() {
	provider := strings.ToLower(strings.TrimSpace(cfg.AWS.Secrets.Provider))
	switch provider {
	case "", "env":
		cfg.AWS.Secrets.Provider = "env"
	case "aws", "secretmanager", "secretsmanager":
		cfg.AWS.Secrets.Provider = "aws"
	default:
		cfg.AWS.Secrets.Provider = "env"
	}
}
