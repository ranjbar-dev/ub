package command

import (
	"context"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"flag"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type retrieveOrdersToRedisCmd struct {
	orderRepository    order.Repository
	redisManager       order.RedisManager
	engineCommunicator order.EngineCommunicator
	logger             platform.Logger
	date               string //the start time to query from
}

func (cmd *retrieveOrdersToRedisCmd) Run(ctx context.Context, flags []string) {
	fmt.Println("start of retrieve open orders command")
	cmd.setNeededData(flags)
	orders := cmd.orderRepository.GetOpenOrders("")
	for _, o := range orders {
		if !o.IsStopOrder() || o.IsSubmitted.Bool == true {
			err := cmd.engineCommunicator.RetrieveOrder(o)
			if err != nil {
				cmd.logger.Error2("error retrieving order to orderbook", err,
					zap.String("service", "retrieveOrdersToRedisCmd"),
					zap.String("method", "Run"),
					zap.Int64("orderID", o.ID),
				)
				continue
			}
			//insert into redis if they do not exist
			continue
		}

		exists, err := cmd.redisManager.Exists(ctx, o)
		if err != nil {
			cmd.logger.Error2("error checking stop order existance in redis", err,
				zap.String("service", "retrieveOrdersToRedisCmd"),
				zap.String("method", "Run"),
				zap.Int64("orderID", o.ID),
			)
			continue
		}
		if !exists {
			err := cmd.redisManager.AddStopOrderToQueue(ctx, o)
			if err != nil {
				cmd.logger.Error2("error adding stop order to queue", err,
					zap.String("service", "retrieveOrdersToRedisCmd"),
					zap.String("method", "Run"),
					zap.Int64("orderID", o.ID),
				)
				continue
			}
		}

	}

	fmt.Println("end of retrieve open orders command")
}

func (cmd *retrieveOrdersToRedisCmd) setNeededData(flags []string) {
	date := flag.String("date", "", "")
	err := flag.CommandLine.Parse(flags)
	if err != nil {
		cmd.logger.Fatal("error in retrieveOrdersToRedisCmd", zap.Error(err))
	}
	t, err := time.Parse("2006-01-02", *date)
	if err == nil {
		startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
		cmd.date = startTime.Format("2006-01-02 15:04:05")
	} else {
		t := time.Now()
		t = t.Add(-1 * 31 * 24 * time.Hour) // 31 day ago
		startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
		cmd.date = startTime.Format("2006-01-02 15:04:05")
	}

}

func NewRetrieveOrderToRedisCmd(orderRepository order.Repository, redisManager order.RedisManager, engineCommunicator order.EngineCommunicator, logger platform.Logger) ConsoleCommand {
	cmd := &retrieveOrdersToRedisCmd{
		orderRepository:    orderRepository,
		redisManager:       redisManager,
		engineCommunicator: engineCommunicator,
		logger:             logger,
	}
	return cmd

}
