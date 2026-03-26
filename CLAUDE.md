# Repository Guidelines

## Project Structure & Module Organization
`cmd/` contains entry points: `cmd/app` builds the CLI, while `cmd/config` and `cmd/gen` support config and model generation. Core application code lives under `internal/` and follows DDD plus Clean Architecture: `domain/` for entities and repository interfaces, `application/usecase/` for orchestration, `interfaces/cli/` for `flag`-based commands, and `infrastructure/` for API, persistence, config, logging, notification, and Wire setup. Runtime config is stored in `configs/config.yaml`.

## Build, Test, and Development Commands
Use `make build` to compile `bin/go-cli-ddd`, and `make run` to build then execute the CLI. Run focused commands such as `make run-account`, `make run-campaign`, or `make run-master` when working on a single command group. Use `make wire` after changing dependency wiring in `internal/infrastructure/wire/`. Test with `make test`, `make test-race`, or `make test-coverage`; run `make test-integration` for tagged integration cases. `make lint` runs the single repo-wide config in `.golangci.yml`, and `make ci` matches the main local pre-PR check.

## Coding Style & Naming Conventions
Target Go `1.25.8` and keep code `gofmt`- and `goimports`-clean. Package names should stay lower-case and concise; exported identifiers use PascalCase, unexported identifiers use camelCase. Logging is standardized on `log/slog`. Keep CLI command definitions in `internal/interfaces/cli/*_command.go`, use case logic in `internal/application/usecase/`, and repository implementations in the matching `internal/infrastructure/` package. Prefer small interfaces in `internal/domain/repository/`.

## Configuration Notes
`prd` is the default base configuration. `dev` overlays `prd`, and `local` overlays `prd` then reads `.env`. If `--env` is omitted, runtime environment is resolved from `ENV`, then `.env`, then `prd`. OS environment variables override YAML in every environment. Switch secret resolution with `SECRETS_PROVIDER=env|aws`: `env` uses direct values from `.env` or the runtime environment, while `aws` uses secret IDs and AWS Secrets Manager.

## Testing Guidelines
This repository uses Go’s `testing` package, with `testify` available for assertions. Place tests beside the code they verify and name them `*_test.go`; reserve `*_integration_test.go` plus the `integration` build tag for external-system or database coverage. Run `make test-race` before opening a PR, and use `make test-coverage` when changing persistence, concurrency, or retry logic.

## Commit & Pull Request Guidelines
Recent history favors short, imperative summaries such as `Added golangci-lint configuration file` or `Improved and more flexible database generation tools`. Keep commits scoped to one change and describe the affected subsystem. PRs should explain the behavioral change, list the validation commands you ran (for example, `make lint && make test-race`), and link the related issue. Include config or schema notes when changes affect `configs/`, persistence, or external API integrations.
