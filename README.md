# Go CLI DDD
![Go CI](https://github.com/yuru-sha/go-cli-ddd/workflows/Go%20CI/badge.svg)


A sample CLI application based on DDD and Clean Architecture using Go 1.24.0, Cobra, GORM, and Google Wire.

## Technology Stack

- Go 1.24.0
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GORM](https://gorm.io/) - ORM library
- [Google Wire](https://github.com/google/wire) - Dependency injection tool
- [zerolog](https://github.com/rs/zerolog) - High-performance structured logging library
- [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) - Goroutine groups with error handling for parallel processing
- [backoff/v4](https://github.com/cenkalti/backoff) - Retry processing with exponential backoff
- [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate) - Rate limiter
- [net/http](https://pkg.go.dev/net/http) - HTTP client
- [golangci-lint](https://golangci-lint.run/) - Linting tool
- [Mermaid](https://mermaid.js.org/) - Text-based diagram creation tool
- [GitHub Actions](https://github.com/features/actions) - CI/CD platform

## Configuration File

The application loads settings from `configs/config.yaml`. The following sections are provided to manage different settings for each environment:

- `local`: Default settings for local development environment
- `dev`: Settings for development environment
- `prd`: Settings for production environment

To specify the environment, use the `--env` flag:

```bash
# Run with default settings (local environment)
./app

# Run with development environment settings
./app --env dev

# Run with production environment settings
./app --env prd
```

Example configuration file:

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

## Project Structure

Project structure following DDD and Clean Architecture principles:

```
.
├── cmd/                    # Application entry points
│   └── app/                # CLI application
├── configs/                # Configuration files
└── internal/               # Private packages
    ├── domain/             # Domain layer
    │   ├── entity/         # Entities
    │   ├── repository/     # Repository interfaces
    │   ├── service/        # Domain services
    │   └── valueobject/    # Value objects
    ├── application/        # Application layer
    │   └── usecase/        # Use cases
    ├── infrastructure/     # Infrastructure layer
    │   ├── config/         # Configuration manager
    │   ├── http/           # HTTP client
    │   ├── logger/         # Logger
    │   ├── persistence/    # Database implementations
    │   ├── notification/   # Notification features
    │   ├── api/            # External API clients
    │   └── wire/           # Dependency injection configuration
    └── interfaces/         # Interface layer
        └── cli/            # CLI interface
```

## Features

This CLI application implements part of an advertising management system and provides the following features:

### Account Management

- Account information synchronization
- Retrieving account information
- Synchronization by specific account ID

### Campaign Management

- Campaign information synchronization
- Retrieving campaign information
- Retrieving campaigns linked to an account ID

## Command Usage Examples

### Account Synchronization

```bash
# Synchronize all accounts
./bin/go-cli-ddd account

# Synchronize only a specific account
./bin/go-cli-ddd account --id 123

# Run with specified sync mode
./bin/go-cli-ddd account --mode diff

# Force synchronization (overwrite existing data)
./bin/go-cli-ddd account --force
```

### Campaign Synchronization

```bash
# Synchronize all campaigns
./bin/go-cli-ddd campaign

# Synchronize only campaigns linked to a specific account
./bin/go-cli-ddd campaign --account-id 123
```

## Setup and Development

### Prerequisites

- Go 1.24.0 or higher
- golangci-lint

### Setup

```bash
# Clone the repository
git clone https://github.com/user/go-cli-ddd.git
cd go-cli-ddd

# Install dependencies
make deps

# Generate dependency injection code
make wire

# Build
make build

# Run
make run
```

### Development Commands

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run tests with Race Detector enabled
make test-race

# Run integration tests
make test-integration

# Run linting
make lint

# Cleanup
make clean

# Run all tasks
make all

# Run CI tasks (lint, test-race, test-coverage, build)
make ci
```

## Continuous Integration

This project implements continuous integration (CI) using GitHub Actions. The following checks are automatically performed:

1. **Lint**: Code quality check using golangci-lint
2. **Test**: Running unit tests (with Race Detector enabled)
3. **Build**: Building the application
4. **Integration**: Running integration tests

The CI workflow is defined in the following file:
- `.github/workflows/ci.yml`

GitHub Actions workflows are automatically executed on pushes to the main branch and pull requests.

## Architecture

This project follows the principles of Domain-Driven Design (DDD) and Clean Architecture:

1. **Domain Layer**: The central layer containing business logic and rules
   - Entities: Business objects (Account, Campaign)
   - Repository Interfaces: Abstraction of data access
   - Domain Services: Operations between entities

2. **Application Layer**: Implementation of use cases
   - Use Cases: Specific functionalities of the application

3. **Infrastructure Layer**: Details of technical implementation
   - Repository Implementations: Concrete implementations of database access
   - Dependency Injection: Management of dependencies between components

4. **Interface Layer**: Interaction with the outside
   - CLI Interface: Interaction with users

## Sequence Diagram

In this project, main processing flows are represented by sequence diagrams using Mermaid notation. This makes it easier to visually understand the execution flow of the code.

### Sequence Diagram for the account Command

Below is a sequence diagram showing the execution flow of the `account` command:

```mermaid
sequenceDiagram
    autonumber

    %% Participant definitions
    actor User
    participant CLI as "CLI<br>(cobra.Command)"
    participant AccountCommand as "AccountCommand<br>(interfaces/cli)"
    participant AccountUseCase as "AccountUseCase<br>(application/usecase)"
    participant AccountAPIRepo as "AccountAPIRepository<br>(infrastructure/api)"
    participant AccountRepo as "AccountRepository<br>(infrastructure/persistence)"
    participant ExternalAPI as "External API<br>(HTTP Server)"
    participant DB as "Database<br>(SQLite)"

    %% Start of sequence
    User->>+CLI: $ go-cli-ddd account [--id=ID] [--mode=MODE] [--force]
    Note over User,CLI: Execute with command line arguments

    %% Delegation from CLI to AccountCommand
    CLI->>+AccountCommand: Execute RunE function
    Note over AccountCommand: Create context.Background
    Note over AccountCommand: Record start time

    %% Flag processing
    Note over AccountCommand: Log flag values

    %% Conditional branching
    alt accountID > 0 (Sync specific account only)
        Note over AccountCommand: Output log for syncing specific account only
        AccountCommand->>+AccountUseCase: SyncAccounts(ctx)
    else accountID == 0 (Sync all accounts)
        AccountCommand->>+AccountUseCase: SyncAccounts(ctx)
    end

    %% Use case processing
    Note over AccountUseCase: Output sync start log

    %% Fetch data from external API
    AccountUseCase->>+AccountAPIRepo: FetchAccounts(ctx)
    AccountAPIRepo->>+ExternalAPI: HTTP GET request
    Note over AccountAPIRepo,ExternalAPI: Mock data is used in actual implementation

    alt Production environment
        ExternalAPI-->>-AccountAPIRepo: Account data (JSON)
    else Development environment
        Note over AccountAPIRepo: Call fetchMockAccounts
        AccountAPIRepo-->>AccountAPIRepo: Generate mock data
    end

    AccountAPIRepo-->>-AccountUseCase: Array of account entities
    Note over AccountUseCase: Log the number of accounts retrieved

    %% Save to database
    AccountUseCase->>+AccountRepo: SaveAll(ctx, accounts)

    %% Transaction processing
    AccountRepo->>+DB: Begin transaction

    loop For each account
        alt Existing account
            AccountRepo->>DB: UPDATE accounts SET ...
        else New account
            AccountRepo->>DB: INSERT INTO accounts ...
        end
    end

    DB-->>-AccountRepo: Save result

    alt If an error occurred
        AccountRepo-->>AccountUseCase: Return error
        AccountUseCase-->>AccountCommand: Return error
        AccountCommand-->>CLI: Return error
        CLI-->>User: Display error message
    else If successful
        AccountRepo-->>-AccountUseCase: Success
        Note over AccountUseCase: Output sync completion log
        AccountUseCase-->>-AccountCommand: Success
        Note over AccountCommand: Calculate elapsed time
        Note over AccountCommand: Output completion log
        AccountCommand-->>-CLI: Success
        CLI-->>-User: Display success message
    end
```

This sequence diagram represents the following processing flow:

1. User executes the `account` command from the command line
2. CLI framework (cobra) calls the RunE function of AccountCommand
3. AccountCommand processes flags and calls the SyncAccounts method of AccountUseCase
4. AccountUseCase retrieves account information from external API (using mock data in development environment)
5. Save the retrieved account information to the database
6. Display the processing result to the user

Sequence diagrams are important documentation for understanding the application's behavior and serve as a reference when adding new features.

## License

This project is released under the [MIT License](LICENSE).

## Contribution

Contributions are welcome! Please feel free to submit pull requests.
