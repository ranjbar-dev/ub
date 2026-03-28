# Task 002: Reduce Handler Boilerplate with Generic Wrapper

## Priority: MEDIUM
## Risk: LOW (refactor only — same behavior, fewer lines)
## Estimated Scope: ~30 handler functions across 8 files

---

## Problem

Almost every handler in `internal/api/handler/` repeats the same pattern:

```go
func SomeHandler(s some.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        p := some.Params{}
        err := c.ShouldBindJSON(&p)
        if err != nil {
            errorResponse, statusCode := HandleValidationError(err)
            c.AbortWithStatusJSON(statusCode, errorResponse)
            return
        }
        u, ok := GetAuthUser(c)
        if !ok {
            return
        }
        resp, statusCode := s.DoSomething(u, p)
        c.JSON(statusCode, resp)
    }
}
```

This 15-line block is duplicated ~30 times. The only differences are:
1. The params struct type
2. The service method called
3. Whether auth is needed
4. Whether extra context (User-Agent, IP) is extracted

## Goal

Create generic wrapper functions that eliminate the boilerplate while keeping the type safety.

## Implementation Plan

### Step 1: Create wrapper functions

Add to `internal/api/handler/common.go`:

```go
// BindAndCall handles the common pattern: bind JSON params → call service → return response.
// Use for public endpoints that don't require authentication.
func BindAndCall[P any](fn func(P) (interface{}, int)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var p P
        if err := c.ShouldBindJSON(&p); err != nil {
            errorResponse, statusCode := HandleValidationError(err)
            c.AbortWithStatusJSON(statusCode, errorResponse)
            return
        }
        resp, statusCode := fn(p)
        c.JSON(statusCode, resp)
    }
}

// AuthBindAndCall handles the common pattern: bind JSON params → get auth user → call service → return response.
// Use for authenticated endpoints.
func AuthBindAndCall[P any](fn func(*user.User, P) (interface{}, int)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var p P
        if err := c.ShouldBindJSON(&p); err != nil {
            errorResponse, statusCode := HandleValidationError(err)
            c.AbortWithStatusJSON(statusCode, errorResponse)
            return
        }
        u, ok := GetAuthUser(c)
        if !ok {
            return
        }
        resp, statusCode := fn(u, p)
        c.JSON(statusCode, resp)
    }
}

// AuthCall handles authenticated endpoints with no JSON body (GET requests).
func AuthCall(fn func(*user.User, *gin.Context) (interface{}, int)) gin.HandlerFunc {
    return func(c *gin.Context) {
        u, ok := GetAuthUser(c)
        if !ok {
            return
        }
        resp, statusCode := fn(u, c)
        c.JSON(statusCode, resp)
    }
}
```

### Step 2: Convert simple handlers

**Before:**
```go
func SetUserProfile(s user.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        p := user.SetUserProfileParams{}
        err := c.ShouldBindJSON(&p)
        if err != nil {
            errorResponse, statusCode := HandleValidationError(err)
            c.AbortWithStatusJSON(statusCode, errorResponse)
            return
        }
        u, ok := GetAuthUser(c)
        if !ok {
            return
        }
        resp, statusCode := s.SetUserProfile(u, p)
        c.JSON(statusCode, resp)
    }
}
```

**After:**
```go
func SetUserProfile(s user.Service) gin.HandlerFunc {
    return AuthBindAndCall(s.SetUserProfile)
}
```

### Step 3: Keep complex handlers as-is

Some handlers need extra logic (e.g., `CreateOrder` extracts User-Agent and builds `UserAgentInfo`). These cannot be simplified with the generic wrapper and should remain manually written. Specifically:

**Do NOT convert these** (they have extra context extraction):
- `handler.CreateOrder` — builds `UserAgentInfo` from User-Agent header
- `handler.Login` — extracts User-Agent and IP into params
- `handler.Register` — same
- `handler.ForgotPassword` — same
- `handler.ForgotPasswordUpdate` — same

### Candidate handlers for conversion

| File | Handler | Wrapper to use |
|------|---------|---------------|
| `handler/user.go` | `SetUserProfile` | `AuthBindAndCall` |
| `handler/user.go` | `Enable2Fa` | `AuthBindAndCall` |
| `handler/user.go` | `Disable2Fa` | `AuthBindAndCall` |
| `handler/user.go` | `ChangePassword` | `AuthBindAndCall` |
| `handler/user.go` | `SendSms` | `AuthBindAndCall` |
| `handler/user.go` | `EnableSms` | `AuthBindAndCall` |
| `handler/user.go` | `DisableSms` | `AuthBindAndCall` |
| `handler/order.go` | `CancelOrder` | `AuthBindAndCall` |
| `handler/payment.go` | `PreWithdraw` | `AuthBindAndCall` |
| `handler/payment.go` | `Withdraw` | `AuthBindAndCall` |
| `handler/payment.go` | `Cancel` | `AuthBindAndCall` |
| `handler/userwithdrawaddress.go` | `NewWithdrawAddress` | `AuthBindAndCall` |
| `handler/userwithdrawaddress.go` | `AddToFavorites` | `AuthBindAndCall` |
| `handler/userwithdrawaddress.go` | `Delete` | `AuthBindAndCall` |
| `handler/currency.go` | `AddOrRemoveFavoritePair` | `AuthBindAndCall` |
| `handler/userbalance.go` | `SetAutoExchange` | `AuthBindAndCall` |
| `handler/user.go` | `GetUserProfile` | `AuthCall` |
| `handler/user.go` | `GetUserData` | `AuthCall` |
| `handler/user.go` | `Get2FaBarcode` | `AuthCall` |
| `handler/userbalance.go` | `PairBalances` | `AuthCall` |
| `handler/userbalance.go` | `AllBalances` | `AuthCall` |

## Verification

```bash
go build ./...
# Behavior is identical — same HTTP status codes, same JSON responses
```

## Notes

- Requires Go 1.18+ for generics (project uses Go 1.22 ✓)
- The `AuthBindAndCall` function requires the service method signature to match `func(*user.User, P) (interface{}, int)`. Verify each candidate's service method matches before converting.
- If a service method returns a concrete type instead of `interface{}`, you may need a thin wrapper lambda.
