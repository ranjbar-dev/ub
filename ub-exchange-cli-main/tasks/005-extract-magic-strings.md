# Task 005: Extract Magic Strings to Constants

## Priority: LOW
## Risk: NONE (constants replace identical literal values)
## Estimated Scope: ~10 files touched

---

## Problem

HTTP header names are used as string literals throughout handler files, repeated 8+ times:

| String | Occurrences | Files |
|--------|-------------|-------|
| `"User-Agent"` | 7+ | auth.go, order.go, configuration.go, user.go |
| `"x-forwarded-for"` | 1 | common.go |
| `"Authorization"` | 2+ | middleware/auth.go |

## Goal

Extract repeated strings into named constants for consistency and discoverability.

## Implementation Plan

### Step 1: Add constants to `internal/api/handler/common.go`

```go
const (
    // HeaderUserAgent is the standard User-Agent HTTP header name.
    HeaderUserAgent = "User-Agent"
    // HeaderXForwardedFor is the proxy-forwarded client IP header.
    HeaderXForwardedFor = "x-forwarded-for"
    // HeaderAuthorization is the standard Authorization HTTP header.
    HeaderAuthorization = "Authorization"
)
```

### Step 2: Replace all occurrences

**`internal/api/handler/auth.go`** (5 occurrences):
```go
// Before:
userAgent := c.GetHeader("User-Agent")
// After:
userAgent := c.GetHeader(HeaderUserAgent)
```

**`internal/api/handler/order.go`** (1 occurrence):
```go
// Before:
userAgentHeader := c.GetHeader("User-Agent")
// After:
userAgentHeader := c.GetHeader(HeaderUserAgent)
```

**`internal/api/handler/common.go`** (1 occurrence in GetClientIP):
```go
// Before:
xForwardedFor := c.GetHeader("x-forwarded-for")
// After:
xForwardedFor := c.GetHeader(HeaderXForwardedFor)
```

**`internal/api/handler/user.go`** — check for User-Agent usage
**`internal/api/handler/configuration.go`** — check for User-Agent usage

**`internal/api/middleware/auth.go`** (Authorization header):
```go
// Before:
tokenString := c.GetHeader("Authorization")
// After:
tokenString := c.GetHeader(handler.HeaderAuthorization)
// OR define the constant in middleware package if circular import
```

### Step 3: Handle circular import risk

If `middleware` imports `handler` and `handler` imports `middleware`, there's a circular dependency. In that case, put the constants in a shared location:
- Option A: Create `internal/api/constants.go` with the constants
- Option B: Put them in the `middleware` package (which handler already imports)

Check existing import graph before deciding:
```bash
grep -rn "\"exchange-go/internal/api/handler\"" internal/api/middleware/
grep -rn "\"exchange-go/internal/api/middleware\"" internal/api/handler/
```

## Verification

```bash
go build ./...
# Verify no raw header strings remain:
grep -rn '"User-Agent"\|"x-forwarded-for"\|"Authorization"' internal/api/ --include="*.go"
# Should return zero results
```
