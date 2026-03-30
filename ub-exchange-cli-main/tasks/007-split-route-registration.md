# Task 007: Split Route Registration into Domain Groups

## Priority: MEDIUM
## Risk: LOW (pure structural refactor — same routes, same middleware)
## Estimated Scope: 1 file split into 1 file + helper functions

---

## Problem

`internal/api/httpserver.go` has a single `registerRoutes()` method (~160 lines, lines 139-296) that registers ALL API routes in one massive function. This makes it hard to:
- Find which routes exist for a specific domain
- Understand which middleware applies to which routes
- Add new routes without scrolling through the entire function

## Current Structure (lines 139-296)

```
registerRoutes()
├── CORS configuration (lines 140-158)
├── Static files (line 160)
├── /api/v1 group
│   ├── /auth (5 routes)
│   ├── /centrifugo (3 routes)
│   ├── /main-data (5 routes)
│   ├── /currencies (8 routes, mixed auth)
│   ├── /withdraw-address (5 routes, auth required)
│   ├── /order (6 routes, auth required)
│   ├── /order-book, /trade-book (2 public routes)
│   ├── /trade (2 routes, auth required)
│   ├── /user-balance (4 routes, auth required)
│   ├── /crypto-payment (5 routes, auth required)
│   ├── /user (11 routes, auth required)
│   └── /user-profile-image (2 routes, auth required)
```

## Goal

Split into domain-specific registration functions while keeping `registerRoutes()` as the orchestrator.

## Implementation Plan

### Step 1: Create `internal/api/routes.go`

Extract each domain group into its own function:

```go
package api

import (
    "exchange-go/internal/api/handler"
    "exchange-go/internal/api/middleware"
    "github.com/gin-gonic/gin"
)

// registerAuthRoutes sets up /auth endpoints (login, register, password reset).
func (s *httpServer) registerAuthRoutes(v1 *gin.RouterGroup) {
    auth := v1.Group("/auth")
    {
        auth.POST("/login", handler.Login(s.services.AuthService))
        auth.POST("/register", handler.Register(s.services.AuthService))
        auth.POST("/forgot-password", handler.ForgotPassword(s.services.AuthService))
        auth.POST("/forgot-password/update", handler.ForgotPasswordUpdate(s.services.AuthService))
        auth.POST("/verify", handler.VerifyEmail(s.services.AuthService))
    }
}

// registerCentrifugoAuthRoutes sets up /centrifugo endpoints for real-time broker authentication.
func (s *httpServer) registerCentrifugoAuthRoutes(v1 *gin.RouterGroup) {
    centrifugoAuth := v1.Group("/centrifugo")
    {
        centrifugoAuth.POST("/login", handler.CentrifugoLogin(s.services.CentrifugoAuthService))
        centrifugoAuth.POST("/acl", handler.CentrifugoACL(s.services.CentrifugoAuthService))
        centrifugoAuth.POST("/superuser", handler.CentrifugoSuperUser(s.services.CentrifugoAuthService))
    }
}

// registerMainDataRoutes sets up /main-data public endpoints (health check, country list, app version).
func (s *httpServer) registerMainDataRoutes(v1 *gin.RouterGroup) {
    mainData := v1.Group("/main-data")
    {
        mainData.GET("/check", handler.Check())
        mainData.GET("/country-list", handler.Countries(s.services.CountryService))
        mainData.GET("/common", handler.GetRecaptchaKey(s.services.ConfigurationService))
        mainData.GET("/version", handler.GetAppVersion(s.services.ConfigurationService))
        mainData.POST("/contact-us", handler.ContactUs(s.services.ConfigurationService))
    }
}

// registerCurrencyRoutes sets up /currencies endpoints (pairs, fees, favorites).
func (s *httpServer) registerCurrencyRoutes(v1 *gin.RouterGroup) {
    // ... move currency routes here
}

// registerOrderRoutes sets up /order endpoints (create, cancel, history).
func (s *httpServer) registerOrderRoutes(v1 *gin.RouterGroup) {
    // ... move order routes here
}

// registerTradeRoutes sets up /trade endpoints (trade history).
func (s *httpServer) registerTradeRoutes(v1 *gin.RouterGroup) {
    // ... move trade routes here
}

// registerUserBalanceRoutes sets up /user-balance endpoints.
func (s *httpServer) registerUserBalanceRoutes(v1 *gin.RouterGroup) {
    // ... move user-balance routes here
}

// registerPaymentRoutes sets up /crypto-payment endpoints (withdraw, deposit).
func (s *httpServer) registerPaymentRoutes(v1 *gin.RouterGroup) {
    // ... move crypto-payment routes here
}

// registerUserRoutes sets up /user endpoints (profile, 2FA, SMS, password).
func (s *httpServer) registerUserRoutes(v1 *gin.RouterGroup) {
    // ... move user routes here
}

// registerWithdrawAddressRoutes sets up /withdraw-address endpoints.
func (s *httpServer) registerWithdrawAddressRoutes(v1 *gin.RouterGroup) {
    // ... move withdraw-address routes here
}

// registerUserProfileImageRoutes sets up /user-profile-image upload/delete endpoints.
func (s *httpServer) registerUserProfileImageRoutes(v1 *gin.RouterGroup) {
    // ... move profile image routes here
}
```

### Step 2: Simplify `registerRoutes()` in httpserver.go

```go
func (s *httpServer) registerRoutes() {
    r := s.engine
    // CORS configuration stays here
    s.configureCORS(r)
    r.Static("/assets", "./assets")

    v1 := r.Group("/api/v1")
    {
        s.registerAuthRoutes(v1)
        s.registerCentrifugoAuthRoutes(v1)
        s.registerMainDataRoutes(v1)
        s.registerCurrencyRoutes(v1)
        s.registerWithdrawAddressRoutes(v1)
        s.registerOrderRoutes(v1)
        v1.GET("/order-book", handler.OrderBook(s.services.OrderBookService))
        v1.GET("/trade-book", handler.TradeBook(s.services.OrderBookService))
        s.registerTradeRoutes(v1)
        s.registerUserBalanceRoutes(v1)
        s.registerPaymentRoutes(v1)
        s.registerUserRoutes(v1)
        s.registerUserProfileImageRoutes(v1)
    }
}
```

## Verification

```bash
go build ./...
# No route changes — same endpoints, same middleware, same handlers
```

## Notes

- The `s.services.AuthService` pattern means these must be methods on `*httpServer`, not standalone functions
- CORS configuration could also be extracted to a `configureCORS(*gin.Engine)` method
- The admin routes (registered in `ListenAndServeAdmin`) should get the same treatment if they have a similar god-function
