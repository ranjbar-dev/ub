package order

import (
	"context"
	"exchange-go/internal/platform"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const UnmatchedOrdersList = "list:unmatched-orders"

// UnmatchedOrdersHandler periodically retries unmatched orders by polling
// a Redis list and resubmitting them to the matching engine.
type UnmatchedOrdersHandler interface {
	// Match runs a polling loop that pops unmatched order IDs from Redis and
	// resubmits them to the matching engine for another matching attempt.
	Match()
}

type unmatchedOrdersHandler struct {
	rc                 platform.RedisClient
	orderRepository    Repository
	engineCommunicator EngineCommunicator
	configs            platform.Configs
	logger             platform.Logger
}

//would start in New func
func (u *unmatchedOrdersHandler) Match() {
	ctx := context.Background()
	var ticker *time.Ticker
	if u.configs.GetEnv() != platform.EnvTest {
		ticker = time.NewTicker(5 * time.Second)
	} else {
		ticker = time.NewTicker(200 * time.Millisecond)
	}
	defer ticker.Stop()
	for range ticker.C {
		orderIDString, err := u.rc.RPop(ctx, UnmatchedOrdersList)
		if err != nil && err != redis.Nil {
			u.logger.Error2("can not get order id from redis", err,
				zap.String("service", "unmatchedOrdersHandler"),
				zap.String("method", "run"),
			)
		}
		if orderIDString == "" {
			continue
		}
		orderID, err := strconv.ParseInt(orderIDString, 10, 64)
		if err != nil {
			u.logger.Error2("can not parse order id", err,
				zap.String("service", "unmatchedOrdersHandler"),
				zap.String("method", "run"),
			)
			continue
		}
		order := &Order{}
		err = u.orderRepository.GetOrderByID(orderID, order)
		if err != nil {
			u.logger.Error2("can not get order by id ", err,
				zap.String("service", "unmatchedOrdersHandler"),
				zap.String("method", "run"),
				zap.Int64("orderID", orderID),
			)
			continue
		}
		if order.Status != StatusOpen {
			u.logger.Warn("order status is not open",
				zap.String("service", "unmatchedOrdersHandler"),
				zap.String("method", "run"),
				zap.Int64("orderID", orderID),
			)
			continue

		}

		err = u.engineCommunicator.SubmitOrder(*order)
		if err != nil {
			u.logger.Error2("can not submit order by id ", err,
				zap.String("service", "unmatchedOrdersHandler"),
				zap.String("method", "run"),
				zap.Int64("orderID", orderID),
			)
		}
		if u.configs.GetEnv() == platform.EnvTest {
			break
		}
	}

}

func NewUnmatchedOrdersHandler(redisClient platform.RedisClient, orderRepository Repository, engineCommunicator EngineCommunicator, configs platform.Configs, logger platform.Logger) UnmatchedOrdersHandler {
	return &unmatchedOrdersHandler{
		rc:                 redisClient,
		orderRepository:    orderRepository,
		engineCommunicator: engineCommunicator,
		configs:            configs,
		logger:             logger,
	}
}
