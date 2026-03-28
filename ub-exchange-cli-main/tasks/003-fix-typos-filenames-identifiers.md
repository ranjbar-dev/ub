# Task 003: Fix Typos in Filenames and Identifiers

## Priority: MEDIUM
## Risk: LOW (rename only — no logic changes)
## Estimated Scope: ~12 files touched

---

## Problem

Several filenames and exported identifiers contain typos that confuse both human readers and AI agents:

| Typo | Correct | Location |
|------|---------|----------|
| `manger` | `manager` | Filenames and identifiers |
| `Unmachtech` | `Unmatched` | Function name |

## Specific Instances

### Typo 1: `ordercreatemanger.go` → `ordercreatemanager.go`

**File:** `internal/order/ordercreatemanger.go`

Rename the file. No identifier changes needed inside (the interface is `OrderCreateManager` which is correct).

### Typo 2: `stopordersubmissionmanger.go` → `stopordersubmissionmanager.go`

**File:** `internal/order/stopordersubmissionmanger.go`

Rename the file. The interface `StopOrderSubmissionManager` is already correct.

### Typo 3: `NewUnmachtechOrdersHandler` → `NewUnmatchedOrdersHandler`

This is a constructor function name that must be updated everywhere it's referenced.

**Definition:**
- `internal/order/unmatchedordershandler.go` — function `NewUnmachtechOrdersHandler`

**References (all must be updated):**
- `internal/di/di_order_services.go` — calls `order.NewUnmachtechOrdersHandler(...)`
- `internal/order/unmatchedordershandler_test.go` — calls `NewUnmachtechOrdersHandler(...)`

## Implementation Plan

### Step 1: Fix filename `ordercreatemanger.go`

```bash
cd internal/order
# Rename file
mv ordercreatemanger.go ordercreatemanager.go
```

Verify no imports reference the filename (Go imports are package-level, not file-level — filenames don't matter for compilation). This is purely for readability.

### Step 2: Fix filename `stopordersubmissionmanger.go`

```bash
cd internal/order
mv stopordersubmissionmanger.go stopordersubmissionmanager.go
```

Same note: filenames don't affect compilation.

### Step 3: Fix `NewUnmachtechOrdersHandler` → `NewUnmatchedOrdersHandler`

**File 1: `internal/order/unmatchedordershandler.go`**
Find and replace the function definition:
```go
// Before:
func NewUnmachtechOrdersHandler(...) UnmatchedOrdersHandler {
// After:
func NewUnmatchedOrdersHandler(...) UnmatchedOrdersHandler {
```

**File 2: `internal/di/di_order_services.go`**
Find and replace the call:
```go
// Before:
return order.NewUnmachtechOrdersHandler(
// After:
return order.NewUnmatchedOrdersHandler(
```

**File 3: `internal/order/unmatchedordershandler_test.go`**
Find and replace all test calls:
```go
// Before:
NewUnmachtechOrdersHandler(
// After:
NewUnmatchedOrdersHandler(
```

### Step 4: Check for any "manger" references in docs

**File: `docs/order-package.md`** contains:
- "redis manger" → "redis manager"
- "stop order submission manger" → "stop order submission manager"

## Verification

```bash
go build ./...
# Verify no old names remain:
grep -rn "Unmachtech\|manger" internal/ docs/ --include="*.go" --include="*.md"
# Should return zero results
```

## Notes

- Go file renames don't break anything since imports are package-level
- The function rename (`NewUnmachtechOrdersHandler` → `NewUnmatchedOrdersHandler`) is the only potentially breaking change — ensure all 3 reference sites are updated atomically
- If this project uses any code generation that references these names, check those configs too
