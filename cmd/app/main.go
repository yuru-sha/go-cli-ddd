package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/wire"
)

func main() {
	// 設定ファイルのパスと環境を取得
	configPath := "configs/config.yaml"
	env := "local"

	// コマンドライン引数から環境を取得
	for i, arg := range os.Args {
		if arg == "--env" && i+1 < len(os.Args) {
			env = os.Args[i+1]
			break
		}
	}

	// アプリケーションの初期化パラメータを作成
	params := wire.AppParams{
		ConfigPath: configPath,
		Env:        env,
	}

	// アプリケーションの初期化
	rootCmd, err := wire.InitializeApp(params)
	if err != nil {
		fmt.Printf("アプリケーションの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// コマンドの実行
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("コマンドの実行に失敗しました")
		os.Exit(1)
	}
}
