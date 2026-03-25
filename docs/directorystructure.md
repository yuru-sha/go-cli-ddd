# ディレクトリ構成

```text
cmd/            エントリーポイント
  app/          本体 CLI
  config/       解決済み設定の確認用
  gen/          GORM Gen 実行用
configs/        環境別 YAML 設定
docs/           補助ドキュメント
internal/
  application/  ユースケース
  domain/       エンティティとリポジトリ契約
  infrastructure/
    api/        外部 API アダプタ
    config/     `prd` 基底 + `.env` / 環境変数上書き
    http/       HTTP クライアントとレスポンス close 共通化
    logger/     `log/slog` 初期化
    notification/
    persistence/
    secrets/
    wire/       DI 構成
  interfaces/
    cli/        Cobra command, request DTO, handler
```
