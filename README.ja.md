# Go CLI DDD
![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)

Go 1.24.0、Cobra、GORM、Google Wireを使用したDDDとクリーンアーキテクチャに基づくCLIアプリケーションのサンプルです。

## 技術スタック

### コア技術
- Go 1.24.0
- DDD（ドメイン駆動設計）
- クリーンアーキテクチャ
- **AIモデル: claude-3-7-sonnet-20250219 (Anthropic Messages API 2023-06-01) ← バージョン変更禁止**

### CLIフレームワーク
- [Cobra](https://github.com/spf13/cobra): ^1.9.1 - CLIフレームワーク
- [Viper](https://github.com/spf13/viper): ^1.19.0 - 設定管理

### データベース
- SQLite: ^1.5.7（GORM SQLiteドライバー）
- [GORM](https://gorm.io/): ^1.25.12 - ORMライブラリ
- GORM Gen: ^0.3.26 - コード生成

### 依存性注入
- [Google Wire](https://github.com/google/wire): ^0.6.0 - 依存性注入ツール

### ロギング
- [zerolog](https://github.com/rs/zerolog): ^1.33.0 - 高性能な構造化ロギングライブラリ

### 並列処理・非同期処理
- [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup): ^0.12.0 - 並列処理のためのエラーハンドリング付きゴルーチングループ
- [backoff/v4](https://github.com/cenkalti/backoff): ^4.3.0 - 指数関数的バックオフによるリトライ処理
- [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate): ^0.11.0 - レートリミッター
- [net/http](https://pkg.go.dev/net/http) - HTTPクライアント

### 開発ツール
- [golangci-lint](https://golangci-lint.run/) - リントツール
- Go Modules - 依存関係管理
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

## AWS Secret Managerの使用方法

このプロジェクトでは、AWS Secret Managerを使用して以下の認証情報を安全に管理することができます：

1. データベース接続情報
2. API認証トークン

### 設定方法

`configs/config.yaml`ファイルで以下の設定を行います：

```yaml
環境名:
  aws:
    region: "ap-northeast-1"  # AWSリージョン
    secrets:
      enabled: true  # Secret Managerを有効にする

  database:
    secret_id: "環境名/database/アプリ名"  # データベース接続情報のシークレットID

  api:
    token_secret_id: "環境名/api/token"  # API認証トークンのシークレットID
```

### シークレットの形式

#### データベース接続情報

```json
{
  "username": "dbuser",
  "password": "dbpassword",
  "host": "db.example.com",
  "port": 3306,
  "dbname": "mydatabase"
}
```

#### API認証トークン

```json
{
  "token": "api-token-value",
  "bearer_token": "bearer-token-value",
  "access_key": "access-key-value",
  "secret_key": "secret-key-value"
}
```

### 認証情報の取得方法

AWS認証情報は環境変数または`~/.aws/credentials`ファイルから自動的に読み込まれます。
詳細は[AWS SDK for Go V2のドキュメント](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/)を参照してください。

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

このCLIアプリケーションは、広告管理システムの一部を実装しており、以下の機能を提供します：

### アカウント管理

- アカウント情報の同期
- アカウント情報の取得
- 特定のアカウントIDによる同期

### キャンペーン管理

- キャンペーン情報の同期
- キャンペーン情報の取得
- アカウントIDに紐づくキャンペーンの取得

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
   - エンティティ: ビジネスオブジェクト（Account, Campaign）
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
    participant CLI as "CLI<br>(cobra v1.9.1)"
    participant Config as "Config<br>(viper v1.19.0)"
    participant AccountCommand as "AccountCommand<br>(interfaces/cli)"
    participant AccountUseCase as "AccountUseCase<br>(application/usecase)"
    participant AccountAPIRepo as "AccountAPIRepository<br>(infrastructure/api)"
    participant HTTPClient as "HTTPClient<br>(net/http)"
    participant RateLimiter as "RateLimiter<br>(rate v0.11.0)"
    participant Backoff as "Backoff<br>(backoff v4.3.0)"
    participant AccountRepo as "AccountRepository<br>(infrastructure/persistence)"
    participant DB as "Database<br>(GORM v1.25.12)"
    participant Logger as "Logger<br>(zerolog v1.33.0)"
    participant 外部API as "外部API<br>(HTTP Server)"

    %% シーケンスの開始
    ユーザー->>+CLI: $ go-cli-ddd account [--id=ID] [--mode=MODE] [--force] [--env=ENV]
    Note over ユーザー,CLI: コマンドライン引数を指定して実行

    %% 設定の読み込み
    CLI->>+Config: 設定ファイルの読み込み (configs/config.yaml)
    Config-->>-CLI: 環境に応じた設定を返却

    %% CLIからAccountCommandへの処理委譲
    CLI->>+AccountCommand: RunE関数を実行
    AccountCommand->>+Logger: ログ初期化 (log_level設定を適用)
    Logger-->>-AccountCommand: ロガーインスタンス

    Note over AccountCommand: context.Backgroundを作成
    Note over AccountCommand: 開始時間を記録

    %% フラグの処理
    Note over AccountCommand: フラグの値をログに記録
    AccountCommand->>+Logger: フラグ情報をログ出力
    Logger-->>-AccountCommand: ログ出力完了

    %% 条件分岐
    alt accountID > 0 (特定のアカウントのみ同期)
        AccountCommand->>+Logger: 特定のアカウントのみ同期するログを出力
        Logger-->>-AccountCommand: ログ出力完了
        AccountCommand->>+AccountUseCase: SyncAccount(ctx, accountID, mode, force)
    else accountID == 0 (全アカウント同期)
        AccountCommand->>+Logger: 全アカウント同期ログを出力
        Logger-->>-AccountCommand: ログ出力完了
        AccountCommand->>+AccountUseCase: SyncAccounts(ctx, mode, force)
    end

    %% ユースケースの処理
    AccountUseCase->>+Logger: 同期開始ログを出力
    Logger-->>-AccountUseCase: ログ出力完了

    %% 外部APIからデータ取得
    AccountUseCase->>+AccountAPIRepo: FetchAccounts(ctx, filter)

    %% HTTPクライアントの設定
    AccountAPIRepo->>+HTTPClient: 新規HTTPクライアント作成
    HTTPClient-->>-AccountAPIRepo: HTTPクライアントインスタンス

    %% レートリミッターの設定
    AccountAPIRepo->>+RateLimiter: NewLimiter(rate, burst)
    RateLimiter-->>-AccountAPIRepo: レートリミッターインスタンス

    %% リトライ設定
    AccountAPIRepo->>+Backoff: NewExponentialBackOff()
    Backoff-->>-AccountAPIRepo: バックオフインスタンス

    %% API呼び出し
    AccountAPIRepo->>+外部API: HTTP GET リクエスト
    Note over AccountAPIRepo,外部API: レートリミットとリトライ処理を適用

    alt 本番環境
        外部API-->>-AccountAPIRepo: アカウントデータ (JSON)
    else 開発環境
        Note over AccountAPIRepo: fetchMockAccountsを呼び出し
        AccountAPIRepo-->>AccountAPIRepo: モックデータを生成
    end

    AccountAPIRepo-->>-AccountUseCase: アカウントエンティティの配列
    AccountUseCase->>+Logger: 取得したアカウント数をログに出力
    Logger-->>-AccountUseCase: ログ出力完了

    %% データベースへの保存
    AccountUseCase->>+AccountRepo: SaveAll(ctx, accounts, mode, force)

    %% トランザクション処理
    AccountRepo->>+DB: トランザクション開始
    DB-->>-AccountRepo: トランザクションインスタンス

    loop 各アカウントについて
        alt mode == "diff" かつ 既存アカウントの場合
            AccountRepo->>DB: 差分のみ更新
        else mode == "full" または 新規アカウントの場合
            alt 既存アカウントの場合
                AccountRepo->>DB: UPDATE accounts SET ...
            else 新規アカウントの場合
                AccountRepo->>DB: INSERT INTO accounts ...
            end
        end
    end

    AccountRepo->>DB: トランザクションコミット
    DB-->>AccountRepo: コミット結果

    alt エラーが発生した場合
        AccountRepo->>DB: トランザクションロールバック
        DB-->>AccountRepo: ロールバック結果
        AccountRepo-->>AccountUseCase: エラー返却
        AccountUseCase->>+Logger: エラーログを出力
        Logger-->>-AccountUseCase: ログ出力完了
        AccountUseCase-->>AccountCommand: エラー返却
        AccountCommand->>+Logger: エラーログを出力
        Logger-->>-AccountCommand: ログ出力完了
        AccountCommand-->>CLI: エラー返却
        CLI-->>ユーザー: エラーメッセージ表示
    else 成功した場合
        AccountRepo-->>-AccountUseCase: 成功 (保存したアカウント数)
        AccountUseCase->>+Logger: 同期完了ログを出力 (処理件数含む)
        Logger-->>-AccountUseCase: ログ出力完了
        AccountUseCase-->>-AccountCommand: 成功
        Note over AccountCommand: 経過時間を計算
        AccountCommand->>+Logger: 完了ログを出力 (経過時間含む)
        Logger-->>-AccountCommand: ログ出力完了
        AccountCommand-->>-CLI: 成功
        CLI-->>-ユーザー: 成功メッセージ表示
    end
```

このシーケンス図は以下の処理フローを表しています：

1. ユーザーがコマンドラインから`account`コマンドを実行し、必要に応じてフラグを指定
2. 設定ファイル（config.yaml）から環境に応じた設定を読み込み
3. CLIフレームワーク（cobra）がAccountCommandのRunE関数を呼び出し
4. ロガー（zerolog）を初期化し、処理開始のログを出力
5. AccountCommandがフラグを処理し、AccountUseCaseのメソッドを呼び出し
6. AccountUseCaseが外部APIからアカウント情報を取得
   - HTTPクライアント、レートリミッター、リトライ処理を適用
   - 開発環境ではモックデータを使用
7. 取得したアカウント情報をデータベース（GORM）に保存
   - トランザクション処理を適用
   - 同期モード（diff/full）に応じた処理
8. 処理結果をログに出力し、ユーザーに表示

シーケンス図は、アプリケーションの動作を理解するための重要なドキュメントであり、新しい機能を追加する際の参考にもなります。

## ライセンス

このプロジェクトは[MITライセンス](LICENSE)の下で公開されています。

## 貢献

貢献を歓迎します！ぜひプルリクエストを送ってください。

## API バージョン管理
### 重要な制約事項
- 外部サービスとの連携は `internal/infrastructure/api/` ディレクトリ内で実装
- これらのファイルは変更禁止（変更が必要な場合は承認が必要）：
  - config.go  - 環境設定の一元管理

### 実装規則
- 環境変数の利用は config.go 経由のみ許可
