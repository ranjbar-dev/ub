package order

import (
	"context"
	"exchange-go/internal/platform"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

const StopOrderQueuePrefix = "queue:stop:order:"

// RedisManager manages stop orders in Redis sorted sets, using the stop-point
// price as the score for efficient range queries when prices move.
type RedisManager interface {
	// AddStopOrderToQueue adds a stop order to the Redis sorted set for its pair and type,
	// using the stop-point price as the score.
	AddStopOrderToQueue(ctx context.Context, o Order) error
	// RemoveStopOrderFromQueue removes a stop order from the Redis sorted set for the given pair.
	RemoveStopOrderFromQueue(ctx context.Context, o Order, pairName string) error
	// GetStopOrdersFromQueue retrieves stop orders whose stop-point prices fall between
	// formerPrice and price, selecting buy and sell queues based on the price movement direction.
	GetStopOrdersFromQueue(ctx context.Context, pairName string, formerPrice string, price string, isPriceRising bool) ([]redis.Z, error)
	// Exists checks whether the given stop order is currently present in the Redis queue.
	Exists(ctx context.Context, o Order) (bool, error)
}

type redisManager struct {
	rc platform.RedisClient
}

func (rm *redisManager) AddStopOrderToQueue(ctx context.Context, o Order) error {
	queue := getQueueName(o.Type, o.Pair.Name)
	stopPointPriceDecimal, err := decimal.NewFromString(o.StopPointPrice.String)
	if err != nil {
		return fmt.Errorf("AddStopOrderToQueue: parse stop point price: %w", err)
	}
	score, _ := stopPointPriceDecimal.Float64()
	member := strconv.FormatInt(o.ID, 10)

	_, err = rm.rc.ZAdd(ctx, queue, score, member)
	if err != nil {
		return fmt.Errorf("AddStopOrderToQueue: zadd: %w", err)
	}
	return nil
}

func (rm *redisManager) RemoveStopOrderFromQueue(ctx context.Context, o Order, pairName string) error {
	queue := getQueueName(o.Type, pairName)
	member := strconv.FormatInt(o.ID, 10)
	_, err := rm.rc.ZRem(ctx, queue, member)
	if err != nil {
		return fmt.Errorf("RemoveStopOrderFromQueue: zrem: %w", err)
	}
	return nil

}

/**
 */
func (rm *redisManager) GetStopOrdersFromQueue(ctx context.Context, pairName string, formerPrice string, price string, isPriceRising bool) ([]redis.Z, error) {
	res := make([]redis.Z, 0)
	buyQueue := getQueueName(TypeBuy, pairName)
	sellQueue := getQueueName(TypeSell, pairName)

	//default params is for price rising
	params := &redis.ZRangeBy{
		Min:    formerPrice,
		Max:    price,
		Offset: 0,
		Count:  10000,
	}

	if !isPriceRising {
		params = &redis.ZRangeBy{
			Min:    price,
			Max:    formerPrice,
			Offset: 0,
			Count:  10000,
		}
	}
	buyOrders, err := rm.rc.ZRangeByScoreWithScores(ctx, buyQueue, params)
	if err != nil {
		return res, fmt.Errorf("GetStopOrdersFromQueue: query buy queue: %w", err)
	}
	res = append(res, buyOrders...)

	sellOrders, err := rm.rc.ZRangeByScoreWithScores(ctx, sellQueue, params)
	if err != nil {
		return res, fmt.Errorf("GetStopOrdersFromQueue: query sell queue: %w", err)
	}
	res = append(res, sellOrders...)

	return res, err
}

func (rm *redisManager) Exists(ctx context.Context, o Order) (bool, error) {
	queue := getQueueName(o.Type, o.Pair.Name)
	member := strconv.FormatInt(o.ID, 10)
	rank, err := rm.rc.ZScore(ctx, queue, member)
	if rank > 0 {
		return true, nil
	}

	if err == redis.Nil {
		return false, nil
	}
	return false, fmt.Errorf("Exists: zscore: %w", err) //it means unknown error
}

func getQueueName(orderType string, pairName string) string {
	return StopOrderQueuePrefix + orderType + ":" + pairName
}

func NewRedisManager(rc platform.RedisClient) RedisManager {
	return &redisManager{rc}
}
