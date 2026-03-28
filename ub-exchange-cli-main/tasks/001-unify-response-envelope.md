# Task 001: Unify API Response Envelope

## Priority: HIGH
## Risk: LOW (no feature breakage — all response formats already contain `status` and `message`)
## Estimated Scope: ~15 files touched

---

## Problem

Three different response structs exist for API error/success responses, causing inconsistency for API consumers and making it harder to add cross-cutting concerns (logging, metrics):

### Current Structs

**1. ErrorResponse** — `internal/api/handler/error.go:9-13`
```go
type ErrorResponse struct {
    Status  bool              `json:"status"`
    Message string            `json:"message"`
    Data    map[string]string `json:"data"`
}
```

**2. Response** — `internal/api/middleware/auth.go:22-25`
```go
type Response struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

**3. InternalServerErrorResponse** — `internal/api/httpserver.go:37-40`
```go
type InternalServerErrorResponse struct {
    Status  bool   `json:"status"`
    Message string `json:"message"`
}
```

## Goal

Replace all three with a single `APIResponse` struct used everywhere.

## Implementation Plan

### Step 1: Create unified response struct

Create `internal/api/handler/response.go`:
```go
package handler

// APIResponse is the standard envelope for all API responses.
// All endpoints must return this format for consistency.
type APIResponse struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse creates an error response with optional validation field errors.
func NewErrorResponse(message string, data map[string]string) APIResponse {
    return APIResponse{
        Status:  false,
        Message: message,
        Data:    data,
    }
}

// NewSuccessResponse creates a success response with optional data payload.
func NewSuccessResponse(message string, data interface{}) APIResponse {
    return APIResponse{
        Status:  true,
        Message: message,
        Data:    data,
    }
}
```

### Step 2: Update handler/error.go

Replace `ErrorResponse` with `APIResponse` in:
- `HandleValidationError()` (line 46) — return `APIResponse` instead of `ErrorResponse`
- `HandleError()` — same treatment

### Step 3: Update middleware/auth.go

Replace `Response{Code: 401, Message: "..."}` with `APIResponse{Status: false, Message: "..."}`.

**WARNING**: This changes the JSON shape from `{"code": 401, "message": "..."}` to `{"status": false, "message": "..."}`. Verify that client apps (ub-app-main, ub-client-cabinet-main) don't depend on the `code` field from auth middleware responses. Search for `"code"` field usage in:
- `ub-app-main/lib/` — Dio interceptors
- `ub-client-cabinet-main/app/` — API error handling

### Step 4: Update httpserver.go

Replace `InternalServerErrorResponse` usage in the recovery middleware (~line 90) with `APIResponse`.

### Step 5: Remove old structs

Delete `ErrorResponse`, `Response`, and `InternalServerErrorResponse` struct definitions.

## Files to Modify

| File | Change |
|------|--------|
| `internal/api/handler/error.go` | Replace `ErrorResponse` → `APIResponse` |
| `internal/api/middleware/auth.go` | Replace `Response` → `APIResponse` |
| `internal/api/httpserver.go` | Replace `InternalServerErrorResponse` → `APIResponse` |
| NEW `internal/api/handler/response.go` | Create unified struct + constructors |

## Verification

```bash
go build ./...
# Then grep to confirm no old types remain:
grep -rn "ErrorResponse\|InternalServerErrorResponse" internal/api/ --include="*.go"
# Should only find the new APIResponse
```

## Caution

The `middleware.Response` struct with `Code int` field is the only breaking shape change. If any client checks the HTTP response body for a `code` field on 401s, that will break. The HTTP status code (401) is unchanged — only the JSON body shape changes.
