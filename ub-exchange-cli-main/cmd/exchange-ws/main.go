package main

import (
	"context"
	"exchange-go/internal/di"
	"exchange-go/internal/externalexchangews"
	"exchange-go/internal/externalexchangews/binance"
	"fmt"
	"os"
	"os/signal"
)

func main() {

	streams := []string{binance.DepthStream}
	if len(os.Args) > 1 {
		streams = os.Args[1:]
	}
	ctx := context.Background()

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	container := di.NewContainer()
	externalExchangeService := container.Get(di.ExternalExchangeWsService).(externalexchangews.Service)
	ws, err := externalExchangeService.GetActiveExternalExchangeWs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
	go func() {
		for {
			select {
			case <-sigChan:
				cancel()
				os.Exit(1)
			}
		}

	}()

	ws.Run(ctx, streams)

}
