package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	OrderBookBidsPrefix = "order-book:bid:"
	OrderBookAsksPrefix = "order-book:ask:"
)

// RedisClient defines Redis sorted set operations specific to order book storage.
type RedisClient interface {
	// ZRangeByScoreWithScores retrieves members from a sorted set within the given score (price) range, including scores.
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error)
	// TxPipeline creates an atomic Redis transaction pipeline for batching multiple commands.
	TxPipeline() redis.Pipeliner
	// ZRem removes a specific member from the sorted set stored at the given key.
	ZRem(ctx context.Context, queue string, member string) (bool, error)
	// ZPopMin removes and returns up to count members with the lowest scores from the sorted set.
	ZPopMin(ctx context.Context, queue string, count int64) ([]redis.Z, error)
	// ZPopMax removes and returns up to count members with the highest scores from the sorted set.
	ZPopMax(ctx context.Context, queue string, count int64) ([]redis.Z, error)
	// ZCount returns the number of members in the sorted set with scores between min and max.
	ZCount(ctx context.Context, key, min, max string) (count int64, err error)
	// ZScore returns the score (price) of a specific member in the sorted set.
	ZScore(ctx context.Context, key string, member string) (float64, error)
}

type redisOrderBookProvider struct {
	rc     RedisClient
	logger Logger
}

func (p *redisOrderBookProvider) GetOrders(ctx context.Context, params OrderBookProviderParams) (orders []Order, err error) {
	key := p.getQueueName(params.Pair, params.Side)
	min := "0"
	max := "+inf"
	if params.Price != "" {
		//limit orders
		if params.Side == SideAsk {
			max = params.Price
		} else {
			min = params.Price
		}
	} else {
		if params.Side == SideAsk {
			max = params.MaxPrice
		} else {
			min = params.MinPrice
		}
	}

	res, err := p.rc.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{Min: min, Max: max, Offset: 0, Count: 10000})
	if err != nil {
		return orders, err
	}

	for _, z := range res {
		o := Order{}
		err := json.Unmarshal([]byte(z.Member.(string)), &o)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (p *redisOrderBookProvider) RewriteOrderBook(ctx context.Context, doneOrders []Order, partialOrder *Order) error {
	if len(doneOrders) == 0 && partialOrder == nil {
		return nil
	}

	pipe := p.rc.TxPipeline()
	for _, order := range doneOrders {
		removingQueueName := p.getQueueName(order.Pair, order.Side)
		res, err := order.MarshalForOrderbook()
		if err != nil {
			return err
		}
		pipe.ZRem(ctx, removingQueueName, string(res))
	}

	if partialOrder != nil {
		partialOrderPrice, err := partialOrder.GetPrice()
		if err != nil {
			return fmt.Errorf("failed to get partial order price: %w", err)
		}
		partialOrderPriceFloat64, exact := partialOrderPrice.Float64()
		if !exact {
			p.logger.Warn("precision loss converting price to float64 for Redis score",
				zap.String("price", partialOrderPrice.String()),
				zap.Float64("float64", partialOrderPriceFloat64),
			)
		}
		res2, err := partialOrder.MarshalForOrderbook()
		if err != nil {
			return fmt.Errorf("failed to marshal partial order: %w", err)
		}
		partialZ := redis.Z{
			Score:  partialOrderPriceFloat64,
			Member: string(res2),
		}
		addedQueueName := p.getQueueName(partialOrder.Pair, partialOrder.Side)
		pipe.ZAdd(ctx, addedQueueName, &partialZ)
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (p *redisOrderBookProvider) getQueueName(pair string, side string) string {
	key := OrderBookBidsPrefix + pair
	if side == SideAsk {
		key = OrderBookAsksPrefix + pair
	}
	return key

}

func (p *redisOrderBookProvider) RemoveOrder(ctx context.Context, order Order) error {
	queue := p.getQueueName(order.Pair, order.Side)
	res, err := order.MarshalForOrderbook()
	if err != nil {
		return err
	}
	_, err = p.rc.ZRem(ctx, queue, string(res))
	return err
}


func (p *redisOrderBookProvider) PopOrders(ctx context.Context, params OrderBookProviderParams) (orders []Order, err error) {
	key := p.getQueueName(params.Pair, params.Side)
	min := "0"
	max := "+inf"
	if params.Price != "" {
		//limit orders
		if params.Side == SideAsk {
			max = params.Price
		} else {
			min = params.Price
		}
	} else {
		if params.Side == SideAsk {
			max = params.MaxPrice
		} else {
			min = params.MinPrice
		}
	}

	// Step 1: fetch all matching members with their scores
	res, err := p.rc.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{Min: min, Max: max, Offset: 0, Count: 10000})
	if err != nil {
		return orders, err
	}

	if len(res) == 0 {
		return orders, nil
	}

	// Step 2: atomically remove exactly the fetched members via TxPipeline
	pipe := p.rc.TxPipeline()
	for _, z := range res {
		pipe.ZRem(ctx, key, z.Member.(string))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return orders, err
	}

	// Step 3: sort results — asks ascending by score, bids descending
	if params.Side == SideAsk {
		sort.Slice(res, func(i, j int) bool {
			return res[i].Score < res[j].Score
		})
	} else {
		sort.Slice(res, func(i, j int) bool {
			return res[i].Score > res[j].Score
		})
	}

	// Step 4: parse and return orders
	for _, z := range res {
		o := Order{}
		err := json.Unmarshal([]byte(z.Member.(string)), &o)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (p *redisOrderBookProvider) Exists(ctx context.Context, order Order) (bool, error) {
	queue := p.getQueueName(order.Pair, order.Side)
	res, err := order.MarshalForOrderbook()
	if err != nil {
		return false, err
	}
	_, err = p.rc.ZScore(ctx, queue, string(res))
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewRedisOrderBookProvider(rc RedisClient, logger Logger) OrderbookProvider {
	return &redisOrderBookProvider{rc: rc, logger: logger}
}
