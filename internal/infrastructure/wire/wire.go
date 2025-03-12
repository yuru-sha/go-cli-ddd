//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/yuru-sha/go-cli-ddd/internal/application/usecase"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/api/externalapi1"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	httpClient "github.com/yuru-sha/go-cli-ddd/internal/infrastructure/http"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/notification"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/persistence/mysql"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/secrets"
	"github.com/yuru-sha/go-cli-ddd/internal/interfaces/cli"
)

// AppParams はアプリケーションのパラメータを表します
type AppParams struct {
	ConfigPath string
	Env        string
}

// InitializeApp はアプリケーションを初期化します
func InitializeApp(params AppParams) (*cobra.Command, error) {
	wire.Build(
		// 設定
		ProvideConfigOptions,
		config.LoadConfig,
		ProvideHTTPConfig,

		// シークレットマネージャー
		secrets.NewAWSSecretsManager,
		ProvideSecretsManager,

		// データベース
		mysql.NewDatabase,
		ProvideDatabaseConnection,
		mysql.NewAccountRepository,
		mysql.NewCampaignRepository,

		// HTTP
		httpClient.NewHTTPClient,

		// ExternalAPI1
		externalapi1.NewAccountRepository,
		externalapi1.NewCampaignRepository,

		// 通知
		notification.NewRepository,

		// ユースケース
		usecase.NewAccountUseCase,
		usecase.NewCampaignUseCase,
		usecase.NewMasterUseCase,

		// コマンド
		cli.NewRootCommand,
		cli.NewAccountCommand,
		cli.NewCampaignCommand,
		cli.NewMasterCommand,

		// ルートコマンドの初期化
		ProvideRootCommand,
	)
	return nil, nil
}

// ProvideConfigOptions は設定オプションを提供します
func ProvideConfigOptions(params AppParams) *config.Options {
	return config.NewConfigOptions(params.ConfigPath, params.Env)
}

// ProvideDatabaseConfig はデータベース設定を提供します
func ProvideDatabaseConfig(cfg *config.Config) *config.DatabaseConfig {
	return &cfg.Database
}

// ProvideHTTPConfig はHTTP設定を提供します
func ProvideHTTPConfig(cfg *config.Config) *config.HTTPConfig {
	return &cfg.HTTP
}

// ProvideAWSConfig はAWS設定を提供します
func ProvideAWSConfig(cfg *config.Config) *config.AWSConfig {
	return &cfg.AWS
}

// ProvideSecretsManager はSecretsManagerインターフェースを提供します
func ProvideSecretsManager(sm *secrets.AWSSecretsManager) secrets.Manager {
	return sm
}

// ProvideDatabaseConnection はデータベース接続を提供します
func ProvideDatabaseConnection(db *mysql.Database) *gorm.DB {
	return db.DB
}

// ProvideRootCommand はルートコマンドを提供します
func ProvideRootCommand(
	rootCmd *cli.RootCommand,
	accountCmd *cli.AccountCommand,
	campaignCmd *cli.CampaignCommand,
	masterCmd *cli.MasterCommand,
) (*cobra.Command, error) {
	rootCmd.Cmd.AddCommand(accountCmd.Cmd)
	rootCmd.Cmd.AddCommand(campaignCmd.Cmd)
	rootCmd.Cmd.AddCommand(masterCmd.Cmd)
	return rootCmd.Cmd, nil
}
