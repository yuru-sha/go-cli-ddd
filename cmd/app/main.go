// Command app starts the main CLI application.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/wire"
	"github.com/yuru-sha/go-cli-ddd/internal/interfaces/cli"
)

func main() {
	opts, _, _, err := cli.ParseRootArgs(os.Args[1:], "configs/config.yaml")
	if err != nil {
		fmt.Printf("起動引数の解析に失敗しました: %v\n", err)
		os.Exit(1)
	}

	params := wire.AppParams{
		ConfigPath: opts.ConfigPath,
		Env:        opts.Env,
	}

	rootCmd, err := wire.InitializeApp(params)
	if err != nil {
		fmt.Printf("アプリケーションの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(os.Args[1:]); err != nil {
		slog.Error("コマンドの実行に失敗しました", "err", err)
		os.Exit(1)
	}
}
