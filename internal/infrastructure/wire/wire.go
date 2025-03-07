//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"

	"github.com/yuru-sha/go-cli-ddd/internal/application/usecase"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/api"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	httpClient "github.com/yuru-sha/go-cli-ddd/internal/infrastructure/http"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/persistence"
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
		ProvideDatabaseConfig,
		ProvideHTTPConfig,
		ProvideAPIConfig,

		// データベース
		persistence.NewDatabase,
		persistence.NewAccountRepository,
		persistence.NewCampaignRepository,

		// HTTP
		httpClient.NewHTTPClient,

		// API
		api.NewAccountAPIRepository,
		api.NewCampaignAPIRepository,

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

// ProvideAPIConfig はAPI設定を提供します
func ProvideAPIConfig(cfg *config.Config) *config.APIConfig {
	return &cfg.API
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
