# Go CLI DDD

Go 1.24.0、Cobra、GORM、Google Wireを使用したDDDとクリーンアーキテクチャに基づくCLIアプリケーションのサンプルです。

![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)

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
- [Mermaid](https://mermaid.js.org/) - テキストベースのダイアグラム作成ツール
- [GitHub Actions](https://github.com/features/actions) - CI/CDプラットフォーム

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

### アカウント管理

- アカウント情報の同期
- アカウント情報の取得
- 特定のアカウントIDによる同期

### キャンペーン管理

- キャンペーン情報の同期
- キャンペーン情報の取得
- アカウントIDに紐づくキャンペーンの取得

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

### アカウント同期

```bash
# 全てのアカウントを同期
./bin/go-cli-ddd account

# 特定のアカウントのみ同期
./bin/go-cli-ddd account --id 123

# 同期モードを指定して実行
./bin/go-cli-ddd account --mode diff

# 強制同期（既存データを上書き）
./bin/go-cli-ddd account --force
```

### キャンペーン同期

```bash
# 全てのキャンペーンを同期
./bin/go-cli-ddd campaign

# 特定のアカウントに紐づくキャンペーンのみ同期
./bin/go-cli-ddd campaign --account-id 123
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

# カバレッジ付きテスト実行
make test-coverage

# Race Detectorを有効にしたテスト実行
make test-race

# 統合テスト実行
make test-integration

# リント実行
make lint

# クリーンアップ
make clean

# 全てのタスク実行
make all

# CI用のタスク実行（lint, test-race, test-coverage, build）
make ci
```

## 継続的インテグレーション

このプロジェクトはGitHub Actionsを使用して継続的インテグレーション（CI）を実装しています。以下のチェックが自動的に実行されます：

1. **Lint**: golangci-lintを使用したコード品質チェック
2. **Test**: ユニットテストの実行（Race Detector有効）
3. **Build**: アプリケーションのビルド
4. **Integration**: 統合テストの実行

CIワークフローは以下のファイルで定義されています：
- `.github/workflows/ci.yml`

GitHub Actionsのワークフローは、mainブランチへのプッシュとプルリクエストで自動的に実行されます。

## アーキテクチャ

このプロジェクトは、ドメイン駆動設計（DDD）とクリーンアーキテクチャの原則に従っています：

1. **ドメイン層**: ビジネスロジックとルールを含む中心的な層
   - エンティティ: ビジネスオブジェクト（Account, Campaign, Task）
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

## シーケンス図

このプロジェクトでは、主要な処理フローをマーメイド記法を使用したシーケンス図で表現しています。これにより、コードの実行フローを視覚的に理解しやすくなります。

### accountコマンドのシーケンス図

以下は、`account`コマンドの実行フローを表すシーケンス図です：

```mermaid
sequenceDiagram
    autonumber

    %% 参加者の定義
    actor ユーザー
    participant CLI as "CLI<br>(cobra.Command)"
    participant AccountCommand as "AccountCommand<br>(interfaces/cli)"
    participant AccountUseCase as "AccountUseCase<br>(application/usecase)"
    participant AccountAPIRepo as "AccountAPIRepository<br>(infrastructure/api)"
    participant AccountRepo as "AccountRepository<br>(infrastructure/persistence)"
    participant 外部API as "外部API<br>(HTTP Server)"
    participant DB as "データベース<br>(SQLite)"

    %% シーケンスの開始
    ユーザー->>+CLI: $ go-cli-ddd account [--id=ID] [--mode=MODE] [--force]
    Note over ユーザー,CLI: コマンドライン引数を指定して実行

    %% CLIからAccountCommandへの処理委譲
    CLI->>+AccountCommand: RunE関数を実行
    Note over AccountCommand: context.Backgroundを作成
    Note over AccountCommand: 開始時間を記録

    %% フラグの処理
    Note over AccountCommand: フラグの値をログに記録

    %% 条件分岐
    alt accountID > 0 (特定のアカウントのみ同期)
        Note over AccountCommand: 特定のアカウントのみ同期するログを出力
        AccountCommand->>+AccountUseCase: SyncAccounts(ctx)
    else accountID == 0 (全アカウント同期)
        AccountCommand->>+AccountUseCase: SyncAccounts(ctx)
    end

    %% ユースケースの処理
    Note over AccountUseCase: 同期開始ログを出力

    %% 外部APIからデータ取得
    AccountUseCase->>+AccountAPIRepo: FetchAccounts(ctx)
    AccountAPIRepo->>+外部API: HTTP GET リクエスト
    Note over AccountAPIRepo,外部API: 実際の実装ではモックデータを使用

    alt 本番環境
        外部API-->>-AccountAPIRepo: アカウントデータ (JSON)
    else 開発環境
        Note over AccountAPIRepo: fetchMockAccountsを呼び出し
        AccountAPIRepo-->>AccountAPIRepo: モックデータを生成
    end

    AccountAPIRepo-->>-AccountUseCase: アカウントエンティティの配列
    Note over AccountUseCase: 取得したアカウント数をログに出力

    %% データベースへの保存
    AccountUseCase->>+AccountRepo: SaveAll(ctx, accounts)

    %% トランザクション処理
    AccountRepo->>+DB: トランザクション開始

    loop 各アカウントについて
        alt 既存アカウントの場合
            AccountRepo->>DB: UPDATE accounts SET ...
        else 新規アカウントの場合
            AccountRepo->>DB: INSERT INTO accounts ...
        end
    end

    DB-->>-AccountRepo: 保存結果

    alt エラーが発生した場合
        AccountRepo-->>AccountUseCase: エラー返却
        AccountUseCase-->>AccountCommand: エラー返却
        AccountCommand-->>CLI: エラー返却
        CLI-->>ユーザー: エラーメッセージ表示
    else 成功した場合
        AccountRepo-->>-AccountUseCase: 成功
        Note over AccountUseCase: 同期完了ログを出力
        AccountUseCase-->>-AccountCommand: 成功
        Note over AccountCommand: 経過時間を計算
        Note over AccountCommand: 完了ログを出力
        AccountCommand-->>-CLI: 成功
        CLI-->>-ユーザー: 成功メッセージ表示
    end
```

このシーケンス図は以下の処理フローを表しています：

1. ユーザーがコマンドラインから`account`コマンドを実行
2. CLIフレームワーク（cobra）がAccountCommandのRunE関数を呼び出し
3. AccountCommandがフラグを処理し、AccountUseCaseのSyncAccountsメソッドを呼び出し
4. AccountUseCaseが外部APIからアカウント情報を取得（開発環境ではモックデータを使用）
5. 取得したアカウント情報をデータベースに保存
6. 処理結果をユーザーに表示

シーケンス図は、アプリケーションの動作を理解するための重要なドキュメントであり、新しい機能を追加する際の参考にもなります。

## ライセンス

MIT
