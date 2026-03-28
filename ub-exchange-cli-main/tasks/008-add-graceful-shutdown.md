# Task 008: Add Graceful Shutdown to exchange-httpd

## Priority: HIGH
## Risk: LOW (additive — only changes shutdown behavior, not request handling)
## Estimated Scope: 1 file

---

## Problem

`cmd/exchange-httpd/main.go` uses `panic()` on errors and has no signal handling:

```go
func main() {
    addr := "0.0.0.0:8000"
    adminAddr := "0.0.0.0:8001"
    if len(os.Args) > 1 {
        addr = os.Args[1]
    }
    container := di.NewContainer()
    httpServer := container.Get(di.HTTPServer).(api.HTTPServer)
    go func() {
        err := httpServer.ListenAndServeAdmin(adminAddr)
        if err != nil {
            panic("can not run admin http server because of err" + err.Error())
        }
    }()
    unmatchedOrdersHandler := container.Get(di.UnmatchedOrderHandler).(order.UnmatchedOrdersHandler)
    go func() {
        unmatchedOrdersHandler.Match()
    }()
    err := httpServer.ListenAndServe(addr)
    if err != nil {
        panic("can not run http server because of err" + err.Error())
    }
}
```

Issues:
1. **No graceful shutdown** — SIGTERM kills the process immediately, dropping in-flight requests
2. **Uses `panic()`** instead of proper error handling and logging
3. **No cleanup** — DI container is never closed (database connections, Redis, RabbitMQ left dangling)
4. **No signal handling** — can't do zero-downtime deploys in Docker/Kubernetes

## Goal

Add signal-based graceful shutdown with a timeout, proper error logging, and DI container cleanup.

## Implementation Plan

### Replace `cmd/exchange-httpd/main.go` with:

```go
package main

import (
    "context"
    "exchange-go/internal/api"
    "exchange-go/internal/di"
    "exchange-go/internal/order"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    addr := "0.0.0.0:8000"
    adminAddr := "0.0.0.0:8001"

    if len(os.Args) > 1 {
        addr = os.Args[1]
    }

    container := di.NewContainer()
    defer container.Delete() // Clean up all DI services on exit

    httpServer := container.Get(di.HTTPServer).(api.HTTPServer)

    // Start admin HTTP server
    go func() {
        if err := httpServer.ListenAndServeAdmin(adminAddr); err != nil {
            fmt.Fprintf(os.Stderr, "admin server failed: %v\n", err)
        }
    }()

    // Start unmatched orders handler
    unmatchedOrdersHandler := container.Get(di.UnmatchedOrderHandler).(order.UnmatchedOrdersHandler)
    go func() {
        unmatchedOrdersHandler.Match()
    }()

    // Start main HTTP server in a goroutine
    errCh := make(chan error, 1)
    go func() {
        if err := httpServer.ListenAndServe(addr); err != nil {
            errCh <- err
        }
    }()

    // Wait for interrupt signal or server error
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    select {
    case sig := <-quit:
        fmt.Fprintf(os.Stdout, "received signal %v, shutting down...\n", sig)
    case err := <-errCh:
        fmt.Fprintf(os.Stderr, "server error: %v\n", err)
        os.Exit(1)
    }

    // Graceful shutdown with 30-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := httpServer.Shutdown(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "graceful shutdown failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("server stopped gracefully")
}
```

### Pre-requisite: Check if HTTPServer interface has Shutdown method

The `api.HTTPServer` interface must expose a `Shutdown(ctx context.Context) error` method. Check `internal/api/httpserver.go` for the interface definition.

**If `Shutdown` doesn't exist**, add it:

```go
// In the HTTPServer interface:
type HTTPServer interface {
    ListenAndServe(addr string) error
    ListenAndServeAdmin(addr string) error
    Shutdown(ctx context.Context) error  // ADD THIS
}

// In the httpServer struct implementation:
func (s *httpServer) Shutdown(ctx context.Context) error {
    // s.engine is *gin.Engine — Gin uses http.Server underneath
    // You need to store the *http.Server when starting
    return s.server.Shutdown(ctx)
}
```

This may require storing the `*http.Server` instance in the `httpServer` struct when `ListenAndServe` is called, so `Shutdown` can reference it.

## Verification

```bash
go build ./cmd/exchange-httpd/
# Run and test signal handling:
./exchange-httpd &
kill -SIGTERM $!
# Should see "received signal terminated, shutting down..." and clean exit
```

## Notes

- The `container.Delete()` call cleans up all DI services (closes DB connections, Redis, RabbitMQ)
- The 30-second timeout gives in-flight requests time to complete before forced shutdown
- This is critical for Docker/Kubernetes deployments where SIGTERM is sent before container stop
- The `unmatchedOrdersHandler.Match()` goroutine should also respect context cancellation — but that's a separate task
