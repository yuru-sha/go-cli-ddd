# Go CLI DDD

Go 1.24.0、Cobra、GORM、Google Wireを使用したDDDとクリーンアーキテクチャに基づく広告管理CLIアプリケーションのサンプルです。

## 技術スタック

- Go 1.24.0
- [Cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [GORM](https://gorm.io/) - ORMライブラリ
- [Google Wire](https://github.com/google/wire) - 依存性注入ツール
- [zerolog](https://github.com/rs/zerolog) - 高性能な構造化ロギングライブラリ
- [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) - 並列処理のためのエラーハンドリング付きゴルーチングループ
- [backoff/v4](https://github.com/cenkalti/backoff) - 指数関数的バックオフによるリトライ処理
- [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate) - レートリミッター
- [net/http](https://pkg.go.dev/net/http) - HTTPクライアント
- [golangci-lint](https://golangci-lint.run/) - リントツール

## 設定ファイル

アプリケーションは`configs/config.yaml`から設定を読み込みます。環境ごとに異なる設定を管理するために、以下のセクションが用意されています：

- `local`: ローカル開発環境用のデフォルト設定
- `dev`: 開発環境用設定
- `prd`: 本番環境用設定

環境を指定するには、`--env`フラグを使用します：

```bash
# デフォルト設定（local環境）で実行
./app

# 開発環境設定で実行
./app --env dev

# 本番環境設定で実行
./app --env prd
```

設定ファイルの例：

```yaml
local:
  app:
    name: "go-cli-ddd"
    debug: true
    log_level: "debug"

  http:
    timeout: 30
    max_retries: 3
    rate_limit:
      qps: 10.0
      burst: 3

  notification:
    slack:
      enabled: false
      webhook_url: "https://hooks.slack.com/services/TXXXXXXXX/BXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
      channel: "#notifications"
      username: "TaskBot"

  external_api:
    task_sync:
      enabled: false
      base_url: "http://localhost:8080"
      api_key: "local_api_key"

dev:
  app:
    debug: true
    log_level: "info"

  notification:
    slack:
      enabled: true
      channel: "#dev-notifications"

  external_api:
    task_sync:
      enabled: true
      base_url: "https://dev-api.example.com"
      api_key: "dev_api_key"

prd:
  app:
    debug: false
    log_level: "info"

  notification:
    slack:
      enabled: true
      channel: "#prd-notifications"
      username: "TaskBot-Production"

  external_api:
    task_sync:
      enabled: true
      base_url: "https://api.example.com"
      api_key: "prd_api_key"
```

## プロジェクト構造

DDDとクリーンアーキテクチャの原則に従ったプロジェクト構造：

```
.
├── cmd/                    # アプリケーションのエントリーポイント
│   └── app/                # CLIアプリケーション
├── configs/                # 設定ファイル
└── internal/               # 非公開パッケージ
    ├── domain/             # ドメイン層
    │   ├── entity/         # エンティティ
    │   ├── repository/     # リポジトリインターフェース
    │   ├── service/        # ドメインサービス
    │   └── valueobject/    # 値オブジェクト
    ├── application/        # アプリケーション層
    │   └── usecase/        # ユースケース
    ├── infrastructure/     # インフラストラクチャ層
    │   ├── config/         # 設定マネージャー
    │   ├── http/           # HTTPクライアント
    │   ├── logger/         # ロガー
    │   ├── parallel/       # 並列処理ユーティリティ
    │   ├── persistence/    # データベース実装
    │   ├── notification/   # 通知機能
    │   ├── api/           # 外部APIクライアント
    │   └── wire/           # 依存性注入設定
    └── interfaces/         # インターフェース層
        └── cli/            # CLIインターフェース
```

## 機能

このCLIアプリケーションは、タスク管理システムを実装しており、以下の機能を提供します：

### タスク管理

- タスクの作成
- タスク一覧の表示（ステータスによるフィルタリング可能）
- タスクステータスの更新
- タスクの削除

### Slack通知

タスクの重要なイベントをSlackに通知します：

- タスク作成時の通知
- タスクステータス更新時の通知
- エラー発生時の通知

### 外部API連携

他システムとのタスク同期をサポートします：

- タスクの外部システムへの同期
- 外部システムからのタスク情報の取得

### タスクの属性

- **ID**: 自動生成されるUUID
- **タイトル**: タスクの名前
- **説明**: タスクの詳細説明
- **ステータス**: TODO、IN_PROGRESS、DONE
- **優先度**: 低、中、高
- **作成日時/更新日時**: タスクの作成/更新日時

## コマンド使用例

### タスク作成

```bash
# 基本的なタスク作成
task create --title "レポート作成" --description "四半期レポートを作成する" --priority 3

# 短縮形
task create -t "メール返信" -d "クライアントからのメールに返信する" -p 2
```

### タスク一覧表示

```bash
# 全てのタスクを表示
task list

# TODOステータスのタスクのみ表示
task list --status TODO

# 進行中のタスクのみ表示
task list --status IN_PROGRESS
```

### タスクステータス更新

```bash
# タスクを進行中に更新
task update-status <タスクID> IN_PROGRESS

# タスクを完了に更新
task update-status <タスクID> DONE
```

### タスク削除

```bash
# タスクを削除
task delete <タスクID>
```

## セットアップと開発

### 前提条件

- Go 1.24.0以上
- golangci-lint

### セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/user/go-cli-ddd.git
cd go-cli-ddd

# 依存関係のインストール
make deps

# 依存性注入コードの生成
make wire

# ビルド
make build

# 実行
make run
```

### 開発コマンド

```bash
# テスト実行
make test

# リント実行
make lint

# クリーンアップ
make clean

# 全てのタスク実行
make all
```

## アーキテクチャ

このプロジェクトは、ドメイン駆動設計（DDD）とクリーンアーキテクチャの原則に従っています：

1. **ドメイン層**: ビジネスロジックとルールを含む中心的な層
   - エンティティ: ビジネスオブジェクト（Task）
   - 値オブジェクト: 不変のプロパティ（TaskID, TaskStatus, TaskPriority）
   - リポジトリインターフェース: データアクセスの抽象化
   - ドメインサービス: エンティティ間の操作

2. **アプリケーション層**: ユースケースの実装
   - ユースケース: アプリケーションの具体的な機能

3. **インフラストラクチャ層**: 技術的な実装の詳細
   - リポジトリ実装: データベースアクセスの具体的な実装
   - 依存性注入: コンポーネント間の依存関係の管理

4. **インターフェース層**: 外部とのインタラクション
   - CLIインターフェース: ユーザーとのインタラクション

## ライセンス

MIT
