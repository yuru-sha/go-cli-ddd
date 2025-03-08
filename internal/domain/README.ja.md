# Domain Layer

## ドメインレイヤーの概要

ドメインレイヤーは、DDDアーキテクチャの中核となる層であり、ビジネスロジックとビジネスルールを表現する場所です。このレイヤーはアプリケーションの「何を行うか」を定義し、技術的な実装の詳細から独立しています。

## 軽量DDDとアンチパターン

伝統的なDDDのディレクトリ構造（以下の「ディレクトリ構造パターン」セクションに示されている）は、軽量DDDの観点からは以下のようなアンチパターンが含まれている可能性があります：

1. **過剰な階層化**: 多くのディレクトリに分割することで、コードの見通しが悪くなり、単純なドメインでも複雑な構造を強制してしまう可能性があります。

2. **形式主義的なパターン適用**: 全てのプロジェクトに対して同じディレクトリ構造を適用することは、プロジェクトの特性や規模に合わない場合があります。

3. **ドメインの断片化**: エンティティ、値オブジェクト、サービスなどを厳格に分離することで、関連するドメイン概念が複数のディレクトリに分散し、全体像を把握しにくくなる場合があります。

4. **過剰な抽象化**: 特に小規模なプロジェクトでは、ファクトリや仕様パターンなどの抽象化が不必要な複雑さを生み出す可能性があります。

## 軽量DDDのアプローチ

軽量DDDでは、以下のようなアプローチが推奨されます：

1. **ドメイン中心の構造化**: 技術的な関心事ではなく、ドメインの境界に基づいてコードを構造化します。

```
domain/
├── account/           # アカウントに関する全てのドメイン要素
│   ├── account.go     # エンティティ、値オブジェクト、ドメインサービス
│   ├── repository.go  # リポジトリインターフェース
│   └── events.go      # ドメインイベント
├── campaign/          # キャンペーンに関する全てのドメイン要素
│   ├── campaign.go
│   ├── repository.go
│   └── service.go
└── shared/            # 共有ドメイン概念
    └── money.go       # 共有値オブジェクト
```

2. **コンテキスト境界の明確化**: 大規模なシステムでは、境界付けられたコンテキスト（Bounded Context）を明示的に分離します。

```
domain/
├── advertising/       # 広告コンテキスト
│   ├── campaign/
│   └── creative/
└── billing/           # 請求コンテキスト
    ├── invoice/
    └── payment/
```

3. **プラグマティックなアプローチ**: プロジェクトの規模や複雑さに応じて、必要なパターンのみを適用します。

## ディレクトリ構造パターン（伝統的なDDD）

以下は伝統的なDDDで一般的に見られるディレクトリ構造です：

```
domain/
├── entity/           # ビジネスエンティティ
├── valueobject/      # 値オブジェクト
├── repository/       # リポジトリインターフェース
├── service/          # ドメインサービス
├── factory/          # ファクトリ（オプション）
├── event/            # ドメインイベント（オプション）
└── specification/    # 仕様パターン（オプション）
```

### entity/

エンティティは、一意の識別子を持ち、ライフサイクルを通じて同一性が保たれるドメインオブジェクトです。

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

// ビジネスロジックを含むメソッド
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

値オブジェクトは、属性のみによって定義され、同一性を持たないイミュータブルなオブジェクトです。

```go
// valueobject/money.go
package valueobject

type Money struct {
    Amount   decimal.Decimal
    Currency string
}

// 値オブジェクトは不変であるべき
func NewMoney(amount decimal.Decimal, currency string) Money {
    return Money{
        Amount:   amount,
        Currency: currency,
    }
}

// 新しいインスタンスを返す操作
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return NewMoney(m.Amount.Add(other.Amount), m.Currency), nil
}
```

### repository/

リポジトリは、エンティティの永続化と取得のための抽象インターフェースを定義します。実装の詳細はインフラストラクチャレイヤーに置かれます。

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

ドメインサービスは、特定のエンティティに自然に属さないビジネスロジックを実装します。

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

// 複数のエンティティを操作するロジック
func (s *AccountService) TransferBetweenAccounts(ctx context.Context, sourceID, targetID int64, amount decimal.Decimal) error {
    // 実装...
}
```

### factory/ (オプション)

ファクトリは、複雑なエンティティや値オブジェクトの生成を担当します。

```go
// factory/campaign_factory.go
package factory

import (
    "github.com/user/app/internal/domain/entity"
    "github.com/user/app/internal/domain/valueobject"
)

type CampaignFactory struct {
    // 依存関係...
}

func NewCampaignFactory() *CampaignFactory {
    return &CampaignFactory{}
}

func (f *CampaignFactory) CreateCampaign(accountID int64, name string, budget valueobject.Money) *entity.Campaign {
    // 複雑な初期化ロジック...
    return &entity.Campaign{
        AccountID: accountID,
        Name:      name,
        Budget:    budget,
        Status:    "draft",
        // その他の初期化...
    }
}
```

### event/ (オプション)

ドメインイベントは、ドメイン内で発生した重要な出来事を表します。

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

// イベントハンドラーインターフェース
type AccountEventHandler interface {
    HandleAccountCreated(event AccountCreated)
    HandleAccountActivated(event AccountActivated)
}
```

### specification/ (オプション)

仕様パターンは、ビジネスルールをカプセル化し、エンティティが特定の条件を満たすかどうかを判断します。

```go
// specification/account_specification.go
package specification

import (
    "github.com/user/app/internal/domain/entity"
)

// 仕様インターフェース
type Specification interface {
    IsSatisfiedBy(entity interface{}) bool
}

// アクティブなアカウントの仕様
type ActiveAccountSpecification struct{}

func (s ActiveAccountSpecification) IsSatisfiedBy(e interface{}) bool {
    account, ok := e.(*entity.Account)
    if !ok {
        return false
    }
    return account.Status == "active"
}
```

## ドメインレイヤーの設計原則

1. **ユビキタス言語の使用**: コードはドメインエキスパートが使用する言語を反映すべきです。
2. **技術的詳細からの独立**: ドメインレイヤーは、データベースやフレームワークなどの技術的な詳細から独立しているべきです。
3. **豊かなドメインモデル**: エンティティや値オブジェクトには、単なるデータ構造ではなく、ビジネスロジックを含めるべきです。
4. **不変条件の強制**: エンティティと値オブジェクトは、常に有効な状態を維持するために不変条件を強制すべきです。
5. **カプセル化**: 内部状態は適切にカプセル化し、ビジネスルールに従った操作のみを許可すべきです。

## 結論：プロジェクトに適したアプローチの選択

DDDの実装方法は、プロジェクトの規模、複雑さ、チームの経験によって異なります。伝統的なDDDの構造は複雑なドメインに適している一方、軽量DDDは小〜中規模のプロジェクトや、DDDを初めて導入するチームに適しています。

最も重要なのは、ドメインの本質を捉え、ビジネスルールを明確に表現することです。構造よりも内容を優先し、プロジェクトの進化に合わせて柔軟にアプローチを調整することが推奨されます。
