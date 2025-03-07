# ディレクトリ構成

以下のディレクトリ構造に従って実装を行ってください：

```
/
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

<!--
### 配置ルール
- UIコンポーネント → `app/components/ui/`
- APIエンドポイント → `app/api/[endpoint]/route.ts`
- 共通処理 → `app/lib/utils/`
- API関連処理 → `src/kumoumi/infrastructure/api/`
-->