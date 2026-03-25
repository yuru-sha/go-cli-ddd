// Command app starts the main CLI application.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/wire"
)

func main() {
	configPath := "configs/config.yaml"
	env := "prd"

	for i, arg := range os.Args {
		if arg == "--env" && i+1 < len(os.Args) {
			env = os.Args[i+1]
		}
	}

	params := wire.AppParams{
		ConfigPath: configPath,
		Env:        env,
	}

	rootCmd, err := wire.InitializeApp(params)
	if err != nil {
		fmt.Printf("アプリケーションの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		slog.Error("コマンドの実行に失敗しました", "err", err)
		os.Exit(1)
	}
}
