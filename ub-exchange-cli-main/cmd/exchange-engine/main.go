package main

import (
	"exchange-go/internal/di"
	"exchange-go/internal/engine"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"os"
	"os/signal"
)

func main() {
	forever := make(chan bool)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill, os.Interrupt)

	container := di.NewContainer()
	rc := container.Get(di.RedisClient).(platform.RedisClient)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := container.Get(di.EngineResultHandler).(order.EngineResultHandler)
	logger := container.Get(di.LoggerService).(platform.Logger)
	e := engine.NewEngine(rc, obp, rh, logger,platform.EnvProd)
	_ = e.SetPostOrderMatchingCall(true)
	e.Run(10,true)
	go func() {
		for {
			select {
			case <-sigChan:
				e.Stop()
				forever <- true
			}
		}
	}()

	<-forever
}
