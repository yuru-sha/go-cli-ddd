# 技術スタック

## コア技術
- Go: ^1.25.8
- DDD（ドメイン駆動設計）
- クリーンアーキテクチャ

## CLI / 設定
- flag（標準ライブラリ）
- os.Getenv（環境変数）
- Viper: ^1.19.0（YAML 設定読み込み）

## データベース
- SQLite: ^1.5.7（GORM SQLiteドライバー）
- GORM: ^1.25.12（ORM）
- GORM Gen: ^0.3.26（コード生成）

## 依存性注入
- Google Wire: ^0.6.0

## ロギング
- log/slog（標準ライブラリ）

## 並列処理・非同期処理
- golang.org/x/sync/errgroup: ^0.12.0
- backoff/v4: ^4.3.0（リトライ処理）
- golang.org/x/time/rate: ^0.11.0（レートリミット）
- net/http（HTTPクライアント）

## 開発ツール
- golangci-lint v2（`.golangci.yml` に一本化）
- Go Modules（依存関係管理）
- Mermaid（テキストベースのダイアグラム作成ツール）
- GitHub Actions（CI/CDプラットフォーム）

## 設定方針
- `prd` を基底設定として使用
- `dev` は `prd` の上書き
- `local` は `prd` の上書き後に `.env` を読み込む
- `--env` 未指定時は `ENV`、次に `.env`、最後に `prd` で決定
- OS 環境変数は全環境で YAML を上書き可能
- `SECRETS_PROVIDER=env|aws` で secrets の取得元を切り替える
