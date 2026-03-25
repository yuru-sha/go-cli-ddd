# Go CLI DDD
![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)

A sample CLI application built with Go `1.25.8`, DDD, Clean Architecture, Cobra, GORM, and Google Wire.

## Technology Stack

- Go `1.25.8`
- Cobra and Viper for CLI and configuration
- GORM and GORM Gen for persistence
- Google Wire for dependency injection
- `log/slog` for structured logging
- `errgroup`, `backoff/v4`, and `rate` for concurrency, retry, and throttling
- `golangci-lint` with a single repo-wide config in [`.golangci.yml`](./.golangci.yml)

## Configuration

Runtime config is defined in [`configs/config.yaml`](./configs/config.yaml).

- `prd` is the base configuration and default runtime environment.
- `dev` overlays `prd`.
- `local` overlays `prd`, then loads `.env`, then applies OS environment variables.
- OS environment variables override YAML for all environments, so ECS task definition variables are effective in `dev` and `prd`.
- Secret Manager is intended for `dev` and `prd`; `.env` is intended for `local`.

Common examples:

```bash
make run
make run-account
./bin/go-cli-ddd --env dev campaign --account-ids 1,2
DATABASE_DSN='file:override.db?cache=shared' ./bin/go-cli-ddd --env prd account
```

Example local overrides:

```dotenv
APP_DEBUG=true
DATABASE_DSN=file:go-cli-ddd.db?cache=shared
EXTERNAL_API1_TOKEN=local-token
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
```

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
  interfaces/   Cobra command boundary
```

## Development

```bash
make install-tools
make wire
make build
make lint
make test
make ci
```

- `make lint` runs the stricter repo-wide lint profile.
- `make wire` regenerates dependency injection glue after constructor changes.
- `make ci` mirrors the main local verification path.

## CLI Design

The CLI layer is split into command definitions, request DTOs, and handlers under [`internal/interfaces/cli`](./internal/interfaces/cli). New commands should follow the same pattern:

1. Define request DTOs and handler interfaces.
2. Keep Cobra flag parsing in the command file.
3. Keep branching and orchestration in the handler/usecase boundary.
4. Register the command through the CLI module set in Wire.
