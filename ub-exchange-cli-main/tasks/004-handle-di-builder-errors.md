# Task 004: Handle DI Builder Registration Errors

## Priority: MEDIUM
## Risk: LOW (adds error visibility — no behavior change on success path)
## Estimated Scope: 6 DI files

---

## Problem

All 109 DI service registrations silently discard errors:

```go
_ = builder.Add(di.Def{
    Name:  cacheService,
    Scope: di.App,
    Build: func(ctn di.Container) (interface{}, error) {
        // ...
    },
})
```

If `builder.Add()` fails (e.g., duplicate name, invalid scope), the error is silently ignored. This makes debugging DI misconfigurations extremely difficult — the app will panic later with a cryptic "service not found" error instead of failing fast at registration time.

## Goal

Replace all `_ = builder.Add(...)` with a helper that panics on registration failure. DI registration is a startup-time operation — failing fast is correct behavior.

## Implementation Plan

### Step 1: Create helper function

Add to `internal/di/container.go`:

```go
// mustAdd registers a service definition with the DI builder.
// Panics if registration fails, ensuring misconfiguration is caught at startup.
func mustAdd(def di.Def) {
    if err := builder.Add(def); err != nil {
        panic(fmt.Sprintf("di: failed to register service %q: %v", def.Name, err))
    }
}
```

Add `"fmt"` to the imports in `container.go`.

### Step 2: Replace all `_ = builder.Add(...)` with `mustAdd(...)`

**Files to update:**
- `internal/di/di_infrastructure.go` — ~13 registrations
- `internal/di/di_repositories.go` — ~25 registrations
- `internal/di/di_services.go` — ~32 registrations
- `internal/di/di_order_services.go` — ~21 registrations
- `internal/di/di_commands.go` — ~16 registrations
- `internal/di/di_http.go` — 1 registration

**Example transformation:**

Before:
```go
func addCacheService() {
    _ = builder.Add(di.Def{
        Name:  cacheService,
        Scope: di.App,
        Build: func(ctn di.Container) (interface{}, error) {
            configsService := ctn.Get(ConfigService).(platform.Configs)
            logger := ctn.Get(LoggerService).(platform.Logger)
            return platform.NewCache(configsService, logger), nil
        },
    })
}
```

After:
```go
func addCacheService() {
    mustAdd(di.Def{
        Name:  cacheService,
        Scope: di.App,
        Build: func(ctn di.Container) (interface{}, error) {
            configsService := ctn.Get(ConfigService).(platform.Configs)
            logger := ctn.Get(LoggerService).(platform.Logger)
            return platform.NewCache(configsService, logger), nil
        },
    })
}
```

### Step 3: Bulk replace

This is a simple find-and-replace across the 6 files:
```
Find:    _ = builder.Add(
Replace: mustAdd(
```

## Verification

```bash
go build ./...
# Count replacements:
grep -rn "mustAdd(" internal/di/ --include="*.go" | wc -l
# Should be ~109
grep -rn "_ = builder.Add(" internal/di/ --include="*.go" | wc -l
# Should be 0
```

## Notes

- Using `panic()` at DI registration time is standard Go practice — the app can't function without its services
- The `builder.Add()` method returns an error only for programming mistakes (duplicate names, nil builder) — never for runtime conditions
- This change makes startup failures immediately visible in logs/Sentry instead of causing cryptic panics later
