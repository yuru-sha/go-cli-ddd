# Go CLI DDD
![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)

Go `1.25.8`、DDD、クリーンアーキテクチャ、`flag`、GORM、Google Wire を使った CLI アプリケーションのサンプルです。

## 技術スタック

- Go `1.25.8`
- `flag` ベースの CLI
- Viper（YAML 設定読み込み）
- GORM / GORM Gen
- Google Wire
- `log/slog`
- リポジトリ共通の lint 設定 [` .golangci.yml`](./.golangci.yml)

## CLI の使い方

メインの実行バイナリは [`cmd/app/main.go`](./cmd/app/main.go) です。

```text
go-cli-ddd [--config path] [--env local|dev|prd] <command> [flags]
```

利用可能なコマンド:

- `account`
- `campaign`
- `master`

実行例:

```bash
make run
make run-account
make run-campaign
make run-master
./bin/go-cli-ddd --env dev campaign --account-ids 1,2 --status active --parallel 3
./bin/go-cli-ddd --env local account --id 1,2 --mode diff --force
./bin/go-cli-ddd --config configs/config.yaml --env prd master --account-ids 1,2 --timeout 30
```

[`internal/interfaces/cli`](./internal/interfaces/cli) の実装時点で使えるフラグ:

- `account`: `--id`, `--mode`, `--force`
- `campaign`: `--account-ids`, `--status`, `--parallel`, `--force`
- `master`: `--account-ids`, `--parallel`, `--timeout`, `--force`

## 設定

実行時設定は [`configs/config.yaml`](./configs/config.yaml) で管理します。

- 既定の実行環境は `prd`
- `dev` は `prd` に対する上書き
- `local` は `prd` に対する上書き後に `.env` を読み込み、その後 OS 環境変数で上書き
- `--env` 未指定時は `ENV`、次に `.env`、最後に `prd` の順で実行環境を決定
- OS 環境変数は全環境で YAML より優先されます
- secrets の解決方式は `SECRETS_PROVIDER=env|aws` で切り替え
- `env` は `.env` または OS 環境変数の直接値を使い、`aws` は設定上の secret ID から AWS Secrets Manager を使います

`.env` の例:

```dotenv
ENV=local
SECRETS_PROVIDER=env
AWS_REGION=ap-northeast-1
APP_DEBUG=true
DATABASE_DSN=file:go-cli-ddd.db?cache=shared
EXTERNAL_API1_TOKEN=local-token
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
```

現在のソースコードで定義されている主な設定領域:

- `app`: アプリ名、debug、log level
- `database`: dialect、DSN、log level、auto migrate、Aurora 設定
- `http`: timeout、retry 回数、rate limit
- `external_api1`: base URL と token または token secret ID
- `external_api2`: base URL と OAuth2 認証情報または secret ID
- `notification.slack`: 有効化、webhook、channel、username、絵文字
- `aws`: region と secrets provider

## ディレクトリ構成

```text
cmd/            エントリーポイント（`app`, `config`, `gen`）
configs/        YAML 設定
docs/           補助ドキュメント
internal/
  application/  ユースケース
  domain/       エンティティとリポジトリ契約
  infrastructure/
    api/        外部サービスアダプタ
    config/     設定読み込みと環境変数上書き
    http/       HTTP ヘルパー
    logger/     slog 初期化
    notification/
    persistence/
    secrets/
    wire/
  interfaces/   `flag` ベースの CLI 境界
scripts/        開発補助スクリプト
```

## 開発コマンド

```bash
make install-tools
make wire
make build
make lint
make test
make test-race
make test-coverage
make test-integration
make watch-go
make ci
```

- `make build` は `bin/go-cli-ddd` を生成します。
- `make run` `make run-account` `make run-campaign` `make run-master` は CLI を実行します。
- `make wire` はコンストラクタ変更後に DI コードを再生成します。
- `make gen-model` は [`cmd/gen/main.go`](./cmd/gen/main.go) を実行します。
- `make lint` は `golangci-lint` を実行します。
- `make watch-go` は [`scripts/watch-go.sh`](./scripts/watch-go.sh) を実行します。
- `make ci` は `lint` `test-race` `test-coverage` `build` をまとめて実行します。

## CLI 設計

CLI 層は [`internal/interfaces/cli`](./internal/interfaces/cli) で、コマンド定義、request DTO、handler に分離しています。新しいコマンドは次の方針で追加します。

1. request DTO と handler interface を定義する
2. `flag.FlagSet` の解釈は command file に閉じ込める
3. 分岐や orchestration は handler / use case 側に寄せる
4. Wire 側の CLI module 登録に追加する
