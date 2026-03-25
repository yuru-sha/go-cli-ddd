package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/logger"
)

var (
	cfgFile string
	env     string
)

// NewRootCommand creates the root Cobra command.
func NewRootCommand() *RootCommand {
	rootCmd := &cobra.Command{
		Use:   "go-cli-ddd",
		Short: "広告管理CLIアプリケーション",
		Long:  `Go 1.25.8、Cobra、GORM、Google Wireを使用したDDDとクリーンアーキテクチャに基づく広告管理CLIアプリケーションです。`,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			cfgOpts := config.NewConfigOptions(cfgFile, env)
			cfg, err := config.LoadConfig(cfgOpts)
			if err != nil {
				return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
			}

			logger.InitLogger(cfg.App.LogLevel, cfg.App.Debug)

			slog.Info(
				"アプリケーションを起動しました",
				"app_name", cfg.App.Name,
				"env", env,
				"log_level", cfg.App.LogLevel,
				"debug", cfg.App.Debug,
			)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.yaml", "設定ファイルのパス")
	rootCmd.PersistentFlags().StringVar(&env, "env", "prd", "実行環境 (local, dev, prd)")

	cobra.OnInitialize(initConfig)

	return &RootCommand{Cmd: rootCmd}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("設定ファイルを使用:", viper.ConfigFileUsed())
	}
}
