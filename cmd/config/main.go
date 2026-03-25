// Command config prints the resolved runtime configuration.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "設定ファイルのパス")
	env := flag.String("env", "prd", "環境（local, dev, prd）")
	flag.Parse()

	opts := config.NewConfigOptions(*configPath, *env)

	cfg, err := config.LoadConfig(opts)
	if err != nil {
		fmt.Printf("設定の読み込みに失敗しました: %v\n", err)
		os.Exit(1)
	}

	jsonBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("JSONへの変換に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("環境: %s (ベース: prd)\n", *env)
	fmt.Println(string(jsonBytes))
}
