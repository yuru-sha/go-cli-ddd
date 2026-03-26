# Go CLI DDD
![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)

A sample CLI application built with Go `1.25.8`, DDD, Clean Architecture, `flag`, GORM, and Google Wire.

## Technology Stack

- Go `1.25.8`
- `flag`-based CLI entrypoint
- Viper for YAML configuration loading
- GORM and GORM Gen for persistence and model generation
- Google Wire for dependency injection
- `log/slog` for structured logging
- `golangci-lint` with a single repo-wide config in [`.golangci.yml`](./.golangci.yml)

## CLI Usage

The main binary is built from [`cmd/app/main.go`](./cmd/app/main.go).

```text
go-cli-ddd [--config path] [--env local|dev|prd] <command> [flags]
```

Available commands:

- `account`
- `campaign`
- `master`

Examples:

```bash
make run
make run-account
make run-campaign
make run-master
./bin/go-cli-ddd --env dev campaign --account-ids 1,2 --status active --parallel 3
./bin/go-cli-ddd --env local account --id 1,2 --mode diff --force
./bin/go-cli-ddd --config configs/config.yaml --env prd master --account-ids 1,2 --timeout 30
```

Command flags implemented in [`internal/interfaces/cli`](./internal/interfaces/cli):

- `account`: `--id`, `--mode`, `--force`
- `campaign`: `--account-ids`, `--status`, `--parallel`, `--force`
- `master`: `--account-ids`, `--parallel`, `--timeout`, `--force`

## Configuration

Runtime configuration is defined in [`configs/config.yaml`](./configs/config.yaml).

- `prd` is the base configuration and the default runtime environment.
- `dev` overlays `prd`.
- `local` overlays `prd`, then values from `.env` are loaded, and finally OS environment variables are applied.
- If `--env` is omitted, the runtime environment is resolved from `ENV`, then `.env`, then defaults to `prd`.
- OS environment variables override YAML values in all environments.
- Secrets resolution is switched by `SECRETS_PROVIDER=env|aws`.
- `env` reads direct values from `.env` or OS environment variables.
- `aws` reads secret IDs from config and resolves them via AWS Secrets Manager.

Example local overrides:

```dotenv
ENV=local
SECRETS_PROVIDER=env
AWS_REGION=ap-northeast-1
APP_DEBUG=true
DATABASE_DSN=file:go-cli-ddd.db?cache=shared
EXTERNAL_API1_TOKEN=local-token
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
```

Key config areas currently present in the source:

- `app`: application name, debug flag, log level
- `database`: dialect, DSN, log level, automigration, Aurora settings
- `http`: timeout, retry count, rate limit
- `external_api1`: base URL and token or token secret ID
- `external_api2`: base URL and OAuth2 credentials or secret ID
- `notification.slack`: enablement, webhook, channel, username, emoji
- `aws`: region and secrets provider

## Project Structure

```text
cmd/            entrypoints (`app`, `config`, `gen`)
configs/        YAML configuration
docs/           supporting design docs
internal/
  application/  use cases
  domain/       entities and repository contracts
  infrastructure/
    api/        external service adapters
    config/     config loading and env overrides
    http/       HTTP client helpers
    logger/     slog setup
    notification/
    persistence/
    secrets/
    wire/
  interfaces/   flag-based command boundary
scripts/        development helper scripts
```

## Development

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

- `make build` builds `bin/go-cli-ddd`.
- `make run`, `make run-account`, `make run-campaign`, and `make run-master` run the CLI entrypoint.
- `make wire` regenerates dependency injection glue after constructor changes.
- `make gen-model` runs [`cmd/gen/main.go`](./cmd/gen/main.go).
- `make lint` runs `golangci-lint`.
- `make watch-go` runs [`scripts/watch-go.sh`](./scripts/watch-go.sh).
- `make ci` runs `lint`, `test-race`, `test-coverage`, and `build`.

## CLI Design

The CLI layer is split into command definitions, request DTOs, and handlers under [`internal/interfaces/cli`](./internal/interfaces/cli). New commands should follow the same pattern:

1. Define request DTOs and handler interfaces.
2. Keep `flag.FlagSet` parsing in the command file.
3. Keep branching and orchestration in the handler/use case boundary.
4. Register the command through the CLI module set in Wire.
