package main

import (
	"context"
	"log"
	"ub-communicator/pkg/di"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	consumer := container.GetConsumer()
	if err := consumer.Consume(context.Background()); err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
}
