package platform

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient provides comprehensive Redis operations for caching, sorted sets
// (used for order books), lists (used for queues), hash maps, and pub/sub messaging.
type RedisClient interface {
	// Get retrieves the string value stored at the given key.
	Get(ctx context.Context, key string) (string, error)
	// Exists reports whether the given key exists in Redis.
	Exists(ctx context.Context, key string) bool
	// HExists reports whether the specified field exists in the hash stored at key.
	HExists(ctx context.Context, key string, field string) bool
	// HGet retrieves the value of a single field in the hash stored at key.
	HGet(ctx context.Context, key string, field string) (string, error)
	// HMGet retrieves the values of multiple fields in the hash stored at key.
	HMGet(ctx context.Context, key string, field ...string) ([]interface{}, error)
	// HGetAll retrieves all fields and values of the hash stored at key.
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	// Set stores a key-value pair in Redis with no expiration.
	Set(ctx context.Context, key string, value interface{}) error
	// HSet sets one or more field-value pairs in the hash stored at key.
	HSet(ctx context.Context, key string, values ...interface{}) error
	// Publish sends a message to the specified Redis pub/sub channel.
	Publish(ctx context.Context, channel string, message interface{}) error
	// TxPipeline returns a pipeliner that wraps queued commands in a Redis transaction (MULTI/EXEC).
	TxPipeline() redis.Pipeliner
	// RPush appends one or more values to the tail of the list stored at key.
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	// RPop removes and returns the last element of the list stored at key.
	RPop(ctx context.Context, key string) (string, error)
	// LPush prepends one or more values to the head of the list stored at key.
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	// LPop removes and returns the first element of the list stored at key.
	LPop(ctx context.Context, key string) (string, error)
	// LRem removes the first count occurrences of value from the list stored at key.
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	// ZRangeByScoreWithScores returns members of the sorted set at key with scores
	// within the specified range, including the scores in the result.
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error)
	// ZAdd adds a member with the given score to the sorted set (e.g., an order book).
	// Returns true if a new member was added, false if the score was updated.
	ZAdd(ctx context.Context, queue string, score float64, member string) (bool, error)
	// ZRem removes a member from the sorted set. Returns true if the member existed.
	ZRem(ctx context.Context, queue string, member string) (bool, error)
	// ZPopMin removes and returns up to count members with the lowest scores from the sorted set.
	ZPopMin(ctx context.Context, queue string, count int64) ([]redis.Z, error)
	// ZPopMax removes and returns up to count members with the highest scores from the sorted set.
	ZPopMax(ctx context.Context, queue string, count int64) ([]redis.Z, error)
	// ZCount returns the number of members in the sorted set at key with scores between min and max.
	ZCount(ctx context.Context, key, min, max string) (int64, error)
	// ZScore returns the score of the specified member in the sorted set at key.
	ZScore(ctx context.Context, key string, member string) (float64, error)
	// Expire sets a time-to-live on the given key. Returns true if the timeout was set.
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	// Del removes the specified keys from Redis. Returns the number of keys that were deleted.
	Del(ctx context.Context, keys ...string) (int64, error)
	// LPos returns the index of the first matching element in the list stored at key.
	LPos(ctx context.Context, key string, value string, a redis.LPosArgs) (int64, error)
	// LRange returns the specified range of elements from the list stored at key.
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// BLPop is a blocking variant of LPop that waits up to timeout for an element
	// to become available on any of the specified keys (used for queue consumers).
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
}
type redisClient struct {
	rc *redis.Client
}

func (rc *redisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := rc.rc.Get(ctx, key).Result()
	return res, err
}

func (rc *redisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	res, err := rc.rc.HGet(ctx, key, field).Result()
	return res, err
}

func (rc *redisClient) HMGet(ctx context.Context, key string, field ...string) ([]interface{}, error) {
	res, err := rc.rc.HMGet(ctx, key, field...).Result()
	return res, err
}

func (rc *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	res, err := rc.rc.HGetAll(ctx, key).Result()
	return res, err
}

func (rc *redisClient) Set(ctx context.Context, key string, value interface{}) error {
	_, err := rc.rc.Set(ctx, key, value, 0).Result()
	return err
}

func (rc *redisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	_, err := rc.rc.HSet(ctx, key, values...).Result()
	return err
}

func (rc *redisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return rc.rc.Publish(ctx, channel, message).Err()

}

func (rc *redisClient) Exists(ctx context.Context, key string) bool {
	res, err := rc.rc.Exists(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	return res > 0
}

func (rc *redisClient) HExists(ctx context.Context, key string, field string) bool {
	res, err := rc.rc.HExists(ctx, key, field).Result()
	if err == redis.Nil {
		return false
	}
	return res
}

func (rc *redisClient) TxPipeline() redis.Pipeliner {
	return rc.rc.TxPipeline()
}

func (rc *redisClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return rc.rc.RPush(ctx, key, values).Result()
}

func (rc *redisClient) RPop(ctx context.Context, key string) (string, error) {
	return rc.rc.RPop(ctx, key).Result()

}

func (rc *redisClient) LPop(ctx context.Context, key string) (string, error) {
	return rc.rc.LPop(ctx, key).Result()

}

func (rc *redisClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return rc.rc.LPush(ctx, key, values...).Result()

}

func (rc *redisClient) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return rc.rc.LRem(ctx, key, count, value).Result()

}

func (rc *redisClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return rc.rc.ZRangeByScoreWithScores(ctx, key, opt).Result()

}

func (rc *redisClient) ZAdd(ctx context.Context, queue string, score float64, member string) (bool, error) {
	z := redis.Z{
		Score:  score,
		Member: member,
	}

	res, err := rc.rc.ZAdd(ctx, queue, &z).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (rc *redisClient) ZRem(ctx context.Context, queue string, member string) (bool, error) {
	res, err := rc.rc.ZRem(ctx, queue, member).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (rc *redisClient) ZPopMin(ctx context.Context, queue string, count int64) ([]redis.Z, error) {
	return rc.rc.ZPopMin(ctx, queue, count).Result()
}

func (rc *redisClient) ZPopMax(ctx context.Context, queue string, count int64) ([]redis.Z, error) {
	return rc.rc.ZPopMax(ctx, queue, count).Result()
}

func (rc *redisClient) ZCount(ctx context.Context, key, min, max string) (count int64, err error) {
	return rc.rc.ZCount(ctx, key, min, max).Result()
}

func (rc *redisClient) ZScore(ctx context.Context, key string, member string) (score float64, err error) {
	return rc.rc.ZScore(ctx, key, member).Result()
}

func (rc *redisClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return rc.rc.Expire(ctx, key, expiration).Result()

}

func (rc *redisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return rc.rc.Del(ctx, keys...).Result()
}

func (rc *redisClient) LPos(ctx context.Context, key string, value string, a redis.LPosArgs) (int64, error) {
	return rc.rc.LPos(ctx, key, value, a).Result()
}

func (rc *redisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rc.rc.LRange(ctx, key, start, stop).Result()
}

func (rc *redisClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return rc.rc.BLPop(ctx, timeout, keys...).Result()
}

func NewRedisClient(c Configs) RedisClient {
	dsn := c.GetString("redis.dsn")
	password := c.GetString("redis.password")
	db := c.GetInt("redis.db")
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       db,
	})
	rc := redisClient{client}
	return &rc
}

func NewRedisTestClient(c *redis.Client) RedisClient {
	return &redisClient{c}

}
