# Go CLI DDD
![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)

Go `1.25.8`、DDD、クリーンアーキテクチャ、`flag`、GORM、Google Wire を使った CLI アプリケーションのサンプルです。

## 技術スタック

- Go `1.25.8`
- `flag` / 環境変数
- Viper（YAML 設定読み込み）
- GORM / GORM Gen
- Google Wire
- `log/slog`
- `errgroup` / `backoff/v4` / `rate`
- リポジトリ共通の lint 設定 [` .golangci.yml`](./.golangci.yml)

## 設定

実行時設定は [`configs/config.yaml`](./configs/config.yaml) で管理します。

- 既定の実行環境は `prd`
- `dev` は `prd` に対する上書き
- `local` は `prd` に対する上書き後に `.env` を読み込み、さらに OS 環境変数で上書き
- `--env` 未指定時は `ENV`、次に `.env`、最後に `prd` の順で実行環境を決定
- OS 環境変数は全環境で有効なので、ECS タスク定義の環境変数でも `dev` / `prd` を上書き可能
- secrets の解決方式は `SECRETS_PROVIDER=env|aws` で切り替え
- `env` は `.env` または OS 環境変数の直接値を使い、`aws` は設定上の secret ID から AWS Secrets Manager を使う

例:

```bash
make run
./bin/go-cli-ddd --env dev campaign --account-ids 1,2
DATABASE_DSN='file:override.db?cache=shared' ./bin/go-cli-ddd --env prd account
```

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
```

## 開発コマンド

```bash
make install-tools
make wire
make build
make lint
make watch-go
make test
make ci
```

- `make lint` は厳しめの共通 lint 設定を実行します。
- `make watch-go` は `*.go` の変更を監視して `gofmt`、`goimports`、`golangci-lint` を再実行します。
- `make wire` はコンストラクタ変更後に DI コードを再生成します。
- `make ci` はローカルでの主要確認手順です。

## CLI 設計

CLI 層は [`internal/interfaces/cli`](./internal/interfaces/cli) で、コマンド定義、request DTO、handler に分離しています。新しいコマンドは次の方針で追加します。

1. request DTO と handler interface を定義する
2. `flag.FlagSet` の解釈は command file に閉じ込める
3. 分岐や orchestration は handler / usecase 側に寄せる
4. Wire 側の CLI module 登録に追加する
