package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

func main() {
	// コマンドライン引数の解析
	configPath := flag.String("config", "configs/config.yaml", "設定ファイルのパス")
	env := flag.String("env", "local", "環境（local, dev, prd）")
	flag.Parse()

	// 設定オプションの作成
	opts := config.NewConfigOptions(*configPath, *env)

	// 設定の読み込み
	cfg, err := config.LoadConfig(opts)
	if err != nil {
		fmt.Printf("設定の読み込みに失敗しました: %v\n", err)
		os.Exit(1)
	}

	// 設定をJSON形式で出力
	jsonBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("JSONへの変換に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("環境: %s (ベース: local)\n", *env)
	fmt.Println(string(jsonBytes))
}
