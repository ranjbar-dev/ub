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
	logger := container.Get(di.LoggerService).(platform.Logger)
	obp := engine.NewRedisOrderBookProvider(rc, logger)
	rh := container.Get(di.EngineResultHandler).(order.EngineResultHandler)
	configs := container.Get(di.ConfigService).(platform.Configs)
	e := engine.NewEngine(rc, obp, rh, logger, configs.GetEnv())
	_ = e.SetPostOrderMatchingCall(true)
	workerCount := configs.GetInt("engine.worker_count")
	if workerCount <= 0 {
		workerCount = 10
	}
	e.Run(workerCount, true)
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
