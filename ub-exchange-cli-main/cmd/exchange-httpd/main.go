package main

import (
	"context"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/order"
	"fmt"
	"net/http"
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
	defer container.Delete()

	httpServer := container.Get(di.HTTPServer).(api.HTTPServer)

	// Start admin HTTP server
	go func() {
		if err := httpServer.ListenAndServeAdmin(adminAddr); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "admin server failed: %v\n", err)
		}
	}()

	// Start unmatched orders handler
	unmatchedOrdersHandler := container.Get(di.UnmatchedOrderHandler).(order.UnmatchedOrdersHandler)
	//TODO should we put this in some other place!??
	go func() {
		unmatchedOrdersHandler.Match()
	}()

	// Start main HTTP server in a goroutine so we can also wait for signals
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
