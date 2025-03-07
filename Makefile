.PHONY: build clean run wire test lint gen-model init-db install-tools test-coverage test-race test-integration ci

# デフォルトターゲット
all: wire build

# ビルド
build:
	go build -o bin/go-cli-ddd ./cmd/app

# クリーン
clean:
	rm -rf bin/
	rm -f go-cli-ddd.db
	# GORM genで生成されたファイルも削除
	rm -rf internal/infrastructure/persistence/query/
	rm -rf internal/infrastructure/persistence/model/

# 実行
run: build
	./bin/go-cli-ddd

# Google Wireによる依存性注入コードの生成
wire:
	cd internal/infrastructure/wire && wire

# GORM genによるモデル生成
gen-model: init-db
	go run ./cmd/gen/main.go

# データベース初期化
init-db:
	mkdir -p internal/infrastructure/persistence/query
	mkdir -p internal/infrastructure/persistence/model

# テスト
test:
	go test -v ./...

# カバレッジ付きテスト
test-coverage:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# Race Detectorを有効にしたテスト
test-race:
	go test -v -race ./...

# 統合テスト
test-integration:
	go test -v -tags=integration ./...

# リント
lint:
	golangci-lint run

# ツールのインストール
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/wire/cmd/wire@latest

# CI用のターゲット
ci: lint test-race test-coverage build

# アカウントコマンドの実行
run-account: build
	./bin/go-cli-ddd account

# キャンペーンコマンドの実行
run-campaign: build
	./bin/go-cli-ddd campaign

# マスターコマンドの実行
run-master: build
	./bin/go-cli-ddd master
