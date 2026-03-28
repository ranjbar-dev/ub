package engine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

const QueueName = "engine:queue:orders"

// QueueHandler provides Redis list operations for the order processing queue.
type QueueHandler interface {
	// RPush appends one or more values to the tail of the list stored at key.
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	// LPush prepends one or more values to the head of the list stored at key.
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	// LPop removes and returns the first element from the list stored at key.
	LPop(ctx context.Context, key string) (string, error)
	// LRem removes occurrences of value from the list stored at key.
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	// LPos returns the position of the first occurrence of value in the list stored at key.
	LPos(ctx context.Context, key string, value string, a redis.LPosArgs) (int64, error)
	// BLPop blocks until an element is available at the head of one of the given lists, or until the timeout expires.
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
}

type queue struct {
	qh QueueHandler
}

func newQueue(qh QueueHandler) *queue {
	return &queue{
		qh: qh,
	}
}

func (q *queue) rPush(ctx context.Context, order Order) error {
	res, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = q.qh.RPush(ctx, QueueName, string(res))
	return err
}

func (q *queue) lPush(ctx context.Context, order Order) error {
	res, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = q.qh.LPush(ctx, QueueName, string(res))
	return err
}

func (q *queue) lPop(ctx context.Context) (Order, error) {
	order := Order{}
	res, err := q.qh.LPop(ctx, QueueName)
	if err != nil {
		return order, err
	}
	err = json.Unmarshal([]byte(res), &order)
	return order, err
}

func (q *queue) blPop(ctx context.Context, timeout time.Duration) (Order, error) {
	order := Order{}
	res, err := q.qh.BLPop(ctx, timeout, QueueName)
	if err != nil {
		return order, err
	}
	err = json.Unmarshal([]byte(res[1]), &order)
	return order, err
}

func (q *queue) remove(ctx context.Context, order Order) error {
	res, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = q.qh.LRem(ctx, QueueName, 0, string(res))
	return err
}

func (q *queue) exists(ctx context.Context, order Order) (bool, error) {
	res, err := json.Marshal(order)
	if err != nil {
		return false, err
	}

	pos, err := q.qh.LPos(ctx, QueueName, string(res), redis.LPosArgs{})
	if err != nil && err != redis.Nil {
		return false, err
	}
	return pos > 0, nil
}
