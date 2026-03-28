# Task 011: Add Documentation to Test Files

## Priority: LOW
## Risk: NONE (comment-only changes)
## Estimated Scope: ~28 test files

---

## Problem

The project has 28 test files with good coverage but zero documentation on:
- What each test suite validates
- How test helpers work
- What seed data is expected
- How to run specific test groups

AI agents working on the codebase cannot quickly determine what is tested, what edge cases are covered, or how to add new tests.

## Goal

Add package-level and function-level godoc to all test files.

## Implementation Plan

### Step 1: Add file-level comments to each test file

**Example format:**
```go
// Tests for the order creation flow. Covers:
// - Limit order creation with valid parameters
// - Market order creation with valid parameters
// - Validation failures (missing pair, invalid side, zero amount)
// - Insufficient balance rejection
// - Duplicate order prevention
//
// Test data: Uses mock repositories with pre-seeded user (ID=1, verified, level=2)
// and currency pairs (BTC/USDT, ETH/USDT).
package order
```

### Step 2: Document test helper functions

Any test helper (`setupTest()`, `newMockRepo()`, `seedData()`, etc.) should have godoc explaining:
- What it sets up
- What defaults it uses
- How to override specific values

### Step 3: Document test data / fixtures

If tests use hardcoded values (user IDs, pair names, amounts), add a comment block at the top of the test file explaining the test data contract:

```go
// Test constants — shared across all tests in this file.
// These must match the mock data returned by the mock repositories.
const (
    testUserID    = 1
    testPairName  = "BTC_USDT"
    testBaseAsset = "BTC"
    // ...
)
```

### Files to document

Find all test files:
```bash
find internal/ -name "*_test.go" -type f
```

Known test files include:
- `internal/order/service_test.go` (1372 lines)
- `internal/order/postordermatchingservice_test.go` (1505 lines)
- `internal/order/ordercreatemanger_test.go`
- `internal/order/decisionmanager_test.go`
- `internal/order/engineresulthandler_test.go`
- `internal/order/eventshandler_test.go`
- `internal/engine/*_test.go`
- `internal/currency/*_test.go`
- Plus ~20 more

## Verification

```bash
go build ./...
go test ./... -count=1
# Tests should still pass — only comments added
```

## Notes

- Start with the two largest test files (1505 and 1372 lines) as they have the most impact
- Use `go doc` to verify the documentation renders correctly
- The mocks package (`internal/mocks/`) should also be documented — explain what each mock implements
