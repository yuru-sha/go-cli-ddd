package cli

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/logger"
)

var (
	cfgFile string
	env     string
)

// NewRootCommand はルートコマンドを作成します
func NewRootCommand() *RootCommand {
	rootCmd := &cobra.Command{
		Use:   "go-cli-ddd",
		Short: "広告管理CLIアプリケーション",
		Long:  `Go 1.24.0、Cobra、GORM、Google Wireを使用したDDDとクリーンアーキテクチャに基づく広告管理CLIアプリケーションです。`,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			// 設定ファイルの読み込み
			cfgOpts := config.NewConfigOptions(cfgFile, env)
			cfg, err := config.LoadConfig(cfgOpts)
			if err != nil {
				return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
			}

			// ロガーの初期化
			logger.InitLogger(cfg.App.LogLevel, cfg.App.Debug)

			log.Info().
				Str("app_name", cfg.App.Name).
				Str("env", env).
				Str("log_level", cfg.App.LogLevel).
				Bool("debug", cfg.App.Debug).
				Msg("アプリケーションを起動しました")

			return nil
		},
	}

	// グローバルフラグの定義
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.yaml", "設定ファイルのパス")
	rootCmd.PersistentFlags().StringVar(&env, "env", "local", "実行環境 (local, dev, prd)")

	// 設定ファイルの読み込み
	cobra.OnInitialize(initConfig)

	return &RootCommand{Cmd: rootCmd}
}

// initConfig は設定ファイルを初期化します
func initConfig() {
	if cfgFile != "" {
		// 指定された設定ファイルを使用
		viper.SetConfigFile(cfgFile)
	} else {
		// デフォルトの設定ファイルを使用
		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 環境変数の読み込み
	viper.AutomaticEnv()

	// 設定ファイルの読み込み
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("設定ファイルを使用:", viper.ConfigFileUsed())
	}
}
