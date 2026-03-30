# Task 012: Create ARCHITECTURE.md with Dependency Graph

## Priority: MEDIUM
## Risk: NONE (documentation only)
## Estimated Scope: 1 new file

---

## Problem

While `AGENTS.md` provides a good package map and data flow overview, there is no dedicated architecture document that shows:
- Service dependency graph (which service depends on which)
- Data flow for key operations (order creation, trade settlement, withdrawal)
- Package layering rules (what can import what)
- Infrastructure topology (Redis sorted sets, RabbitMQ queues, Centrifugo channels)

AI agents need this to understand impact analysis when modifying code.

## Goal

Create `ARCHITECTURE.md` at project root with visual dependency graphs and data flow diagrams.

## Content Outline

### Section 1: Package Dependency Layers

```
┌─────────────────────────────────────────────┐
│                  cmd/                        │  Entry points (4 binaries)
├─────────────────────────────────────────────┤
│              internal/api/                   │  HTTP handlers, middleware, routing
├─────────────────────────────────────────────┤
│  internal/order/  internal/auth/  ...        │  Domain services (business logic)
├─────────────────────────────────────────────┤
│          internal/repository/                │  Data access (GORM)
├─────────────────────────────────────────────┤
│           internal/platform/                 │  Infrastructure abstractions
├─────────────────────────────────────────────┤
│     MySQL    Redis    RabbitMQ    Centrifugo       │  External systems
└─────────────────────────────────────────────┘
```

### Section 2: Key Data Flows

Document these critical paths:
1. **Order Creation Flow**: Client → handler → CreateOrder → DecisionManager → Engine → Redis queue
2. **Order Matching Flow**: Engine worker → Redis order book → EngineResultHandler → PostOrderMatchingService → DB + Centrifugo
3. **Trade Settlement Flow**: PostOrderMatchingService → updateBalances → createTransactions → pushDataToUsers
4. **Withdrawal Flow**: handler → PaymentService → WalletService → external blockchain API
5. **Real-time Data Flow**: Trade event → LiveDataService → Centrifugo server → Client WebSocket

### Section 3: DI Service Graph

Document the ~110 services and their dependency relationships. Group by domain:
- Infrastructure services (config, cache, DB, Redis, RabbitMQ, Centrifugo, logger)
- Repository services (20+ repositories)
- Domain services (auth, user, order, currency, payment, etc.)
- API services (HTTP server, handlers)
- CLI commands (16 commands)

### Section 4: Redis Data Structures

Document what's stored in Redis and the key patterns:
- Order books: sorted sets with `pair:side` keys
- Cached data: `entity:id` pattern with TTLs
- In-queue orders: lists
- Live ticker data: hash maps

### Section 5: RabbitMQ Queue Topology

Document queues, exchanges, routing keys, and which service produces/consumes each.

### Section 6: Centrifugo Channel Structure

Document Centrifugo channels used for real-time push to clients (tickers, order updates, trade notifications).

## Implementation Plan

1. Read `AGENTS.md` for existing documentation
2. Read `internal/di/container.go` and split files for service dependency info
3. Trace key flows through the codebase
4. Create `ARCHITECTURE.md` with Mermaid diagrams where possible

## Verification

- Document should be accurate against current codebase
- No code changes needed
