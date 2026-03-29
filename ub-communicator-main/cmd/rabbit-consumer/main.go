package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ub-communicator/pkg/di"
)

const shutdownTimeout = 15 * time.Second

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer container.GetLogger().Shutdown(2 * time.Second)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		consumer := container.GetConsumer()
		errChan <- consumer.Consume(ctx)
	}()

	select {
	case sig := <-sigChan:
		log.Printf("Received signal %v, initiating graceful shutdown...", sig)
		cancel()
		select {
		case <-errChan:
			log.Println("Consumer stopped gracefully")
		case <-time.After(shutdownTimeout):
			log.Println("Shutdown timeout exceeded, forcing exit")
		}
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Consumer failed: %v", err)
		}
	}
}
