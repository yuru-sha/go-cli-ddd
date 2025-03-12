# 技術スタック

## コア技術
- Go: ^1.24.0
- DDD（ドメイン駆動設計）
- クリーンアーキテクチャ
- **AIモデル: claude-3-7-sonnet-20250219 (Anthropic Messages API 2023-06-01) ← バージョン変更禁止**

## CLIフレームワーク
- Cobra: ^1.9.1
- Viper: ^1.19.0（設定管理）

## データベース
- SQLite: ^1.5.7（GORM SQLiteドライバー）
- GORM: ^1.25.12（ORM）
- GORM Gen: ^0.3.26（コード生成）

## 依存性注入
- Google Wire: ^0.6.0

## ロギング
- zerolog: ^1.33.0

## 並列処理・非同期処理
- golang.org/x/sync/errgroup: ^0.12.0
- backoff/v4: ^4.3.0（リトライ処理）
- golang.org/x/time/rate: ^0.11.0（レートリミット）
- net/http（HTTPクライアント）

## 開発ツール
- golangci-lint（リントツール）
- Go Modules（依存関係管理）
- Mermaid（テキストベースのダイアグラム作成ツール）
- GitHub Actions（CI/CDプラットフォーム）

---

# API バージョン管理
## 重要な制約事項
- 外部サービスとの連携は `internal/infrastructure/api/` ディレクトリ内で実装
- これらのファイルは変更禁止（変更が必要な場合は承認が必要）：
  - config.go  - 環境設定の一元管理

## 実装規則
- 環境変数の利用は config.go 経由のみ許可
