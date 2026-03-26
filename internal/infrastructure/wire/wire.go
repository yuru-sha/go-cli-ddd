//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
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

type AppParams struct {
	ConfigPath string
	Env        string
}

func InitializeApp(params AppParams) (*cli.RootCommand, error) {
	wire.Build(
		ProvideConfigOptions,
		config.LoadConfig,
		ProvideHTTPConfig,
		secrets.NewManager,
		mysql.NewDatabase,
		ProvideDatabaseConnection,
		mysql.NewAccountRepository,
		mysql.NewCampaignRepository,
		httpClient.NewHTTPClient,
		externalapi1.NewAccountRepository,
		externalapi1.NewCampaignRepository,
		notification.NewRepository,
		usecase.NewAccountUseCase,
		usecase.NewCampaignUseCase,
		usecase.NewMasterUseCase,
		cli.NewAccountHandler,
		cli.NewCampaignHandler,
		cli.NewMasterHandler,
		cli.NewAccountCommand,
		cli.NewCampaignCommand,
		cli.NewMasterCommand,
		ProvideCommandModules,
		ProvideRootCommand,
	)
	return nil, nil
}

func ProvideConfigOptions(params AppParams) *config.Options {
	return config.NewConfigOptions(params.ConfigPath, params.Env)
}

func ProvideDatabaseConfig(cfg *config.Config) *config.DatabaseConfig {
	return &cfg.Database
}

func ProvideHTTPConfig(cfg *config.Config) *config.HTTPConfig {
	return &cfg.HTTP
}

func ProvideAWSConfig(cfg *config.Config) *config.AWSConfig {
	return &cfg.AWS
}

func ProvideDatabaseConnection(db *mysql.Database) *gorm.DB {
	return db.DB
}

func ProvideCommandModules(
	accountCmd *cli.AccountCommand,
	campaignCmd *cli.CampaignCommand,
	masterCmd *cli.MasterCommand,
) []cli.CommandModule {
	return []cli.CommandModule{accountCmd, campaignCmd, masterCmd}
}

func ProvideRootCommand(modules []cli.CommandModule) (*cli.RootCommand, error) {
	rootCmd := cli.NewRootCommand()
	for _, module := range modules {
		rootCmd.Register(module)
	}
	return rootCmd, nil
}
