# Domain Layer

## Overview of the Domain Layer

The domain layer is the core layer of DDD architecture and is where business logic and business rules are expressed. This layer defines "what the application does" and is independent of the technical implementation details.

## Lightweight DDD and Anti-patterns

The traditional DDD directory structure (shown in the "Directory Structure Pattern" section below) may contain the following anti-patterns from a lightweight DDD perspective:

1. **Excessive layering**: Dividing code into many directories can make the code less readable and may force complex structures even for simple domains.

2. **Formalistic pattern application**: Applying the same directory structure to all projects may not be suitable for the characteristics and scale of each project.

3. **Domain fragmentation**: Strictly separating entities, value objects, services, etc. can scatter related domain concepts across multiple directories, making it difficult to grasp the overall picture.

4. **Excessive abstraction**: Especially in small-scale projects, abstractions like factories and specification patterns may create unnecessary complexity.

## Lightweight DDD Approach

In lightweight DDD, the following approaches are recommended:

1. **Domain-centered structuring**: Structure code based on domain boundaries rather than technical concerns.

```
domain/
├── account/           # All domain elements related to accounts
│   ├── account.go     # Entities, value objects, domain services
│   ├── repository.go  # Repository interfaces
│   └── events.go      # Domain events
├── campaign/          # All domain elements related to campaigns
│   ├── campaign.go
│   ├── repository.go
│   └── service.go
└── shared/            # Shared domain concepts
    └── money.go       # Shared value objects
```

2. **Clarification of context boundaries**: In large-scale systems, explicitly separate bounded contexts.

```
domain/
├── advertising/       # Advertising context
│   ├── campaign/
│   └── creative/
└── billing/           # Billing context
    ├── invoice/
    └── payment/
```

3. **Pragmatic approach**: Apply only the necessary patterns according to the scale and complexity of the project.

## Directory Structure Pattern (Traditional DDD)

Below is a directory structure commonly seen in traditional DDD:

```
domain/
├── entity/           # Business entities
├── valueobject/      # Value objects
├── repository/       # Repository interfaces
├── service/          # Domain services
├── factory/          # Factories (optional)
├── event/            # Domain events (optional)
└── specification/    # Specification patterns (optional)
```

### entity/

Entities are domain objects that have a unique identifier and maintain their identity throughout their lifecycle.

```go
// entity/account.go
package entity

type Account struct {
    ID        int64
    Name      string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Method containing business logic
func (a *Account) Activate() error {
    if a.Status == "active" {
        return errors.New("account is already active")
    }
    a.Status = "active"
    a.UpdatedAt = time.Now()
    return nil
}
```

### valueobject/

Value objects are immutable objects defined only by their attributes and do not have identity.

```go
// valueobject/money.go
package valueobject

type Money struct {
    Amount   decimal.Decimal
    Currency string
}

// Value objects should be immutable
func NewMoney(amount decimal.Decimal, currency string) Money {
    return Money{
        Amount:   amount,
        Currency: currency,
    }
}

// Operations return new instances
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return NewMoney(m.Amount.Add(other.Amount), m.Currency), nil
}
```

### repository/

Repositories define abstract interfaces for persisting and retrieving entities. Implementation details are placed in the infrastructure layer.

```go
// repository/account_repository.go
package repository

import (
    "context"

    "github.com/user/app/internal/domain/entity"
)

type AccountRepository interface {
    FindByID(ctx context.Context, id int64) (*entity.Account, error)
    Save(ctx context.Context, account *entity.Account) error
    Delete(ctx context.Context, id int64) error
    FindAll(ctx context.Context) ([]*entity.Account, error)
}
```

### service/

Domain services implement business logic that does not naturally belong to a specific entity.

```go
// service/account_service.go
package service

import (
    "context"

    "github.com/user/app/internal/domain/entity"
    "github.com/user/app/internal/domain/repository"
)

type AccountService struct {
    accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
    return &AccountService{
        accountRepo: accountRepo,
    }
}

// Logic that operates on multiple entities
func (s *AccountService) TransferBetweenAccounts(ctx context.Context, sourceID, targetID int64, amount decimal.Decimal) error {
    // Implementation...
}
```

### factory/ (Optional)

Factories are responsible for creating complex entities or value objects.

```go
// factory/campaign_factory.go
package factory

import (
    "github.com/user/app/internal/domain/entity"
    "github.com/user/app/internal/domain/valueobject"
)

type CampaignFactory struct {
    // Dependencies...
}

func NewCampaignFactory() *CampaignFactory {
    return &CampaignFactory{}
}

func (f *CampaignFactory) CreateCampaign(accountID int64, name string, budget valueobject.Money) *entity.Campaign {
    // Complex initialization logic...
    return &entity.Campaign{
        AccountID: accountID,
        Name:      name,
        Budget:    budget,
        Status:    "draft",
        // Other initializations...
    }
}
```

### event/ (Optional)

Domain events represent important occurrences within the domain.

```go
// event/account_events.go
package event

import (
    "time"

    "github.com/user/app/internal/domain/entity"
)

type AccountCreated struct {
    Account   *entity.Account
    Timestamp time.Time
}

type AccountActivated struct {
    AccountID int64
    Timestamp time.Time
}

// Event handler interface
type AccountEventHandler interface {
    HandleAccountCreated(event AccountCreated)
    HandleAccountActivated(event AccountActivated)
}
```

### specification/ (Optional)

The specification pattern encapsulates business rules and determines whether an entity meets specific conditions.

```go
// specification/account_specification.go
package specification

import (
    "github.com/user/app/internal/domain/entity"
)

// Specification interface
type Specification interface {
    IsSatisfiedBy(entity interface{}) bool
}

// Active account specification
type ActiveAccountSpecification struct{}

func (s ActiveAccountSpecification) IsSatisfiedBy(e interface{}) bool {
    account, ok := e.(*entity.Account)
    if !ok {
        return false
    }
    return account.Status == "active"
}
```

## Domain Layer Design Principles

1. **Use of ubiquitous language**: Code should reflect the language used by domain experts.
2. **Independence from technical details**: The domain layer should be independent of technical details such as databases and frameworks.
3. **Rich domain model**: Entities and value objects should include business logic, not just be data structures.
4. **Enforcement of invariants**: Entities and value objects should enforce invariants to maintain valid states at all times.
5. **Encapsulation**: Internal state should be properly encapsulated, allowing only operations that comply with business rules.

## Conclusion: Choosing the Right Approach for Your Project

The implementation method of DDD varies depending on the scale, complexity, and team experience of the project. While traditional DDD structures are suitable for complex domains, lightweight DDD is more appropriate for small to medium-sized projects or teams introducing DDD for the first time.

The most important thing is to capture the essence of the domain and clearly express business rules. It is recommended to prioritize content over structure and flexibly adjust your approach as the project evolves.
