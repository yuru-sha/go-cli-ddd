package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config はアプリケーション設定を管理します
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	API      APIConfig      `mapstructure:"api"`
}

// AppConfig はアプリケーション全般の設定です
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

// DatabaseConfig はデータベース接続の設定です
type DatabaseConfig struct {
	Dialect     string `mapstructure:"dialect"`
	DSN         string `mapstructure:"dsn"`
	LogLevel    string `mapstructure:"log_level"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
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

// APIConfig は外部APIの設定です
type APIConfig struct {
	Account  APIEndpointConfig `mapstructure:"account"`
	Campaign APIEndpointConfig `mapstructure:"campaign"`
}

// APIEndpointConfig はAPIエンドポイントの設定です
type APIEndpointConfig struct {
	BaseURL  string `mapstructure:"base_url"`
	Endpoint string `mapstructure:"endpoint"`
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

	// 指定された環境の設定を取得
	var config Config
	if err := viper.UnmarshalKey(strings.ToLower(env), &config); err != nil {
		return nil, fmt.Errorf("設定の解析に失敗しました: %w", err)
	}

	return &config, nil
}
