# Task 013: Document DI Service Registration Order and Dependencies

## Priority: MEDIUM
## Risk: NONE (comment-only changes)
## Estimated Scope: 7 DI files

---

## Problem

The DI container registers ~110 services in a specific order in `NewContainer()`. The order matters because `ctn.Get()` inside `Build` functions resolves services that must already be registered. However:
- There are no comments explaining why services are registered in this order
- There's no documentation of which service depends on which
- Adding a new service requires manually tracing dependencies to find the correct insertion point

## Current Registration Order (container.go lines 127-235)

The `NewContainer()` function calls ~110 `addXxx()` functions in a specific order. The order is implicitly dependency-driven but not documented.

## Goal

Add structured comments to `NewContainer()` and header comments to each DI file explaining the dependency chain.

## Implementation Plan

### Step 1: Add section comments in container.go `NewContainer()`

```go
func NewContainer() di.Container {
    if builder == nil {
        builder, _ = di.NewBuilder(di.App)
    }

    // === Infrastructure (no dependencies) ===
    addConfigService()
    addCacheService()        // depends on: configService, loggerService
    addDBClient()            // depends on: configService
    addLogger()              // depends on: configService
    addWSClient()            // depends on: configService, loggerService
    addCentrifugoClient()          // depends on: configService, loggerService
    addRedisClient()         // depends on: configService, loggerService
    addHTTPClient()          // depends on: configService, loggerService

    // === Messaging ===
    addRabbitmqClient()      // depends on: configService, loggerService
    addQueueManager()        // depends on: rabbitmqClient, loggerService
    addCentrifugoManager()         // depends on: centrifugoClient, loggerService

    // === Repositories (depend on: dbClient) ===
    addOrderRepository()
    addUserBalanceRepository()
    addCurrencyRepository()
    // ... etc

    // === Domain Services (depend on: repositories + infrastructure) ===
    addCurrencyService()     // depends on: currencyRepo, pairRepo, cacheService
    addUserService()         // depends on: userRepo, passwordEncoder, ...
    // ... etc

    // === Order Engine (depend on: domain services) ===
    addEngine()
    addEngineCommunicator()
    // ... etc

    // === CLI Commands (depend on: domain services) ===
    addDeleteCacheCommand()
    // ... etc

    // === HTTP Server (depends on: everything) ===
    addHTTPServer()

    return builder.Build()
}
```

### Step 2: Add header comments to each DI file

**di_infrastructure.go:**
```go
// DI registrations for infrastructure services.
// These are registered first as they have no internal dependencies.
// Other services depend on: configService, loggerService, cacheService, dbClient, redisClient.
```

**di_repositories.go:**
```go
// DI registrations for data access repositories.
// All repositories depend on dbClient (GORM).
// Repositories are stateless — they can be registered in any order relative to each other.
```

**di_services.go:**
```go
// DI registrations for domain services.
// Services depend on repositories and infrastructure.
// Registration order matters: services used by other services must be registered first.
```

**di_order_services.go:**
```go
// DI registrations for order/trading domain services.
// These form the core trading engine pipeline:
// DecisionManager → Engine → EngineCommunicator → EngineResultHandler → PostOrderMatchingService
```

### Step 3: Document the dependency chain for critical paths

Add a comment block at the top of `container.go`:
```go
// Service Dependency Chain (critical path — order matching):
//
//   configService → all services
//   dbClient → all repositories
//   redisClient → orderRedisManager, engine, orderbookService
//   currencyRepository → currencyService → orderService, paymentService
//   userRepository → userService → authService
//   orderRepository → orderService, engineResultHandler
//   decisionManager → orderCreateManager
//   engine → engineCommunicator → orderService
//   postOrderMatchingService → engineResultHandler
//   engineResultHandler → (consumes Redis queue, calls postOrderMatchingService)
```

## Verification

```bash
go build ./...
# Comments only — no logic changes
```
