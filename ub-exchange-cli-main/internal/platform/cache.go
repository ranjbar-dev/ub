package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const CachePrefix = "ec:" //exchange cache

// Cache provides a high-level application cache with TTL and lazy-loading support,
// backed by go-redis/cache. All keys are automatically prefixed with CachePrefix.
// Cache operations are disabled in the test environment.
type Cache interface {
	// Get retrieves the cached value for the given key and unmarshals it into value.
	// Returns cache.ErrCacheMiss if the key does not exist.
	Get(ctx context.Context, key string, value interface{}) error
	// Set stores a value in the cache with the specified TTL. The optional do function
	// enables lazy-loading: it is called to compute the value on a cache miss.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration, do func(*cache.Item) (interface{}, error)) error
	// Delete removes one or more keys from the cache.
	Delete(ctx context.Context, keys ...string) error
	// DeleteAll removes all cached keys matching the given glob pattern.
	// An empty pattern defaults to "*", clearing the entire cache namespace.
	DeleteAll(ctx context.Context, pattern string) error
}

type cacheClient struct {
	rc      *redis.Client
	cache   *cache.Cache
	configs Configs
	logger  Logger
}

func (c *cacheClient) Get(ctx context.Context, key string, value interface{}) error {
	if c.configs.GetEnv() == EnvTest { //to be sure no cache  is used in test env
		return fmt.Errorf("test env")
	}
	key = CachePrefix + key
	err := c.cache.Get(ctx, key, value)
	//	err := cache.cache.GetSkippingLocalCache(ctx, key, value)
	if err != nil {
		if err != cache.ErrCacheMiss {
			c.logger.Warn("error in cache get",
				zap.Error(err),
				zap.String("service", "cache"),
				zap.String("method", "Get"),
			)
		}
	}
	return err
}

func (c *cacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration, do func(*cache.Item) (interface{}, error)) error {
	if c.configs.GetEnv() == EnvTest {
		return nil
	}
	key = CachePrefix + key
	item := &cache.Item{
		Ctx:            ctx,
		Key:            key,
		Value:          value,
		TTL:            ttl,
		Do:             do,
		SkipLocalCache: true,
	}
	err := c.cache.Set(item)
	if err != nil {
		c.logger.Warn("error in cache set",
			zap.Error(err),
			zap.String("service", "cache"),
			zap.String("method", "Set"),
		)
	}
	return err
}

func (c *cacheClient) Delete(ctx context.Context, keys ...string) error {
	var prefixedKeys []string
	for _, key := range keys {
		prefixedKeys = append(prefixedKeys, CachePrefix+key)
	}
	return c.delete(ctx, prefixedKeys...)
}

func (c *cacheClient) delete(ctx context.Context, keys ...string) error {
	if len(keys) > 0 {
		_, err := c.rc.Del(ctx, keys...).Result()
		if err != nil {
			c.logger.Warn("error in cache delele",
				zap.Error(err),
				zap.String("service", "cache"),
				zap.String("method", "Set"),
			)
		}
	}

	return nil
}

func (c *cacheClient) DeleteAll(ctx context.Context, pattern string) error {
	//TODO should we close the client or not at the end
	if pattern == "" {
		pattern = "*"
	}
	cachePattern := CachePrefix + pattern
	iter := c.rc.Scan(ctx, 0, cachePattern, 0).Iterator()
	for iter.Next(ctx) {
		c.rc.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func NewCache(configs Configs, logger Logger) Cache {
	dsn := configs.GetString("redis.dsn")
	password := configs.GetString("redis.password")
	db := configs.GetInt("redis.db")
	rc := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       db,
	})

	c := cache.New(&cache.Options{
		Redis:      rc,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &cacheClient{
		rc:      rc,
		cache:   c,
		configs: configs,
		logger:  logger,
	}
}
