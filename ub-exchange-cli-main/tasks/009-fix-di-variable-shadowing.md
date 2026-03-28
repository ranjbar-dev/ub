# Task 009: Fix Variable Shadowing in DI Registration Functions

## Priority: LOW
## Risk: NONE (variable rename only — no logic change)
## Estimated Scope: 6 DI files, ~80+ instances

---

## Problem

Throughout the DI registration files, local variables shadow their same-named constants:

```go
const userService = "userService"  // package-level constant

func addSomeService() {
    _ = builder.Add(di.Def{
        Build: func(ctn di.Container) (interface{}, error) {
            userService := ctn.Get(userService).(user.Service)  // local shadows constant!
            // ...
        },
    })
}
```

This compiles because Go allows shadowing, but it's confusing — the `userService` on the left is a local variable, while the `userService` inside `ctn.Get()` is the package constant. If someone accidentally uses the local variable where they meant the constant (or vice versa), it silently does the wrong thing.

## Instances (80+ across 6 files)

Examples:
- `di_commands.go:72` — `userService := ctn.Get(userService).(user.Service)`
- `di_commands.go:74` — `currencyService := ctn.Get(currencyService).(currency.Service)`
- `di_services.go:43` — `liveDataService := ctn.Get(liveDataService).(livedata.Service)`
- `di_services.go:94` — `userService := ctn.Get(userService).(user.Service)`
- `di_order_services.go:109` — `currencyService := ctn.Get(currencyService).(currency.Service)`

This pattern repeats for nearly every DI registration.

## Goal

Rename local variables to avoid shadowing. Use a consistent prefix convention.

## Implementation Plan

### Naming Convention

Add a descriptive suffix to local variables:

```go
// Before (shadowing):
userService := ctn.Get(userService).(user.Service)

// After (clear):
userSvc := ctn.Get(userService).(user.Service)
```

**Convention:**
| Constant pattern | Local variable pattern |
|---|---|
| `xxxService` | `xxxSvc` |
| `xxxRepository` | `xxxRepo` |
| `xxxManager` | `xxxMgr` |
| `xxxClient` | `xxxCli` |
| `xxxHandler` | `xxxHdl` |

### Files to modify

| File | Approximate instances |
|------|----------------------|
| `internal/di/di_infrastructure.go` | ~13 |
| `internal/di/di_repositories.go` | ~15 |
| `internal/di/di_services.go` | ~40 |
| `internal/di/di_order_services.go` | ~25 |
| `internal/di/di_commands.go` | ~16 |
| `internal/di/di_http.go` | ~5 |

### Example Transformation

**Before (`di_commands.go`):**
```go
func addUbUpdateUserWalletBalancesCommand() {
    _ = builder.Add(di.Def{
        Name:  UpdateUserWalletBalancesCommand,
        Build: func(ctn di.Container) (interface{}, error) {
            userService := ctn.Get(userService).(user.Service)
            userBalanceService := ctn.Get(userBalanceService).(userbalance.Service)
            currencyService := ctn.Get(currencyService).(currency.Service)
            logger := ctn.Get(LoggerService).(platform.Logger)
            return command.NewUbUpdateUserWalletBalances(userService, userBalanceService, currencyService, logger), nil
        },
    })
}
```

**After:**
```go
func addUbUpdateUserWalletBalancesCommand() {
    _ = builder.Add(di.Def{
        Name:  UpdateUserWalletBalancesCommand,
        Build: func(ctn di.Container) (interface{}, error) {
            userSvc := ctn.Get(userService).(user.Service)
            userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
            currencySvc := ctn.Get(currencyService).(currency.Service)
            loggerSvc := ctn.Get(LoggerService).(platform.Logger)
            return command.NewUbUpdateUserWalletBalances(userSvc, userBalanceSvc, currencySvc, loggerSvc), nil
        },
    })
}
```

## Verification

```bash
go build ./...
# Optionally run go vet to confirm no shadowing warnings:
go vet ./internal/di/...
```

## Notes

- This is a large but purely mechanical change — perfect for automated find/replace
- Each function is self-contained (local variables don't escape the closure), so changes are isolated
- Can be combined with Task 004 (DI builder error handling) for a single pass through the DI files
