package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache wraps Redis client with caching functionality
type RedisCache struct {
	client *redis.Client
}

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(config Config) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisCache{
		client: rdb,
	}
}

// NewRedisCacheFromURL creates a new Redis cache from URL
func NewRedisCacheFromURL(url string) (*RedisCache, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	rdb := redis.NewClient(opts)
	return &RedisCache{client: rdb}, nil
}

// Ping tests the Redis connection
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Set stores a value in Redis with TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from Redis
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

// Delete removes a key from Redis
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeletePattern deletes all keys matching a pattern
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in Redis
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return count > 0, nil
}

// SetNX sets a key only if it doesn't exist (for locking)
func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.SetNX(ctx, key, data, ttl).Result()
}

// Increment increments a counter
func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// IncrementWithTTL increments a counter and sets TTL if it's a new key
func (c *RedisCache) IncrementWithTTL(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	pipe := c.client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment with TTL: %w", err)
	}
	
	return incrCmd.Val(), nil
}

// GetMultiple retrieves multiple values from Redis
func (c *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	values, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple values: %w", err)
	}

	result := make(map[string]interface{})
	for i, key := range keys {
		if values[i] != nil {
			var value interface{}
			if err := json.Unmarshal([]byte(values[i].(string)), &value); err == nil {
				result[key] = value
			}
		}
	}

	return result, nil
}

// SetMultiple stores multiple values in Redis
func (c *RedisCache) SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	pipe := c.client.Pipeline()
	
	for key, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		pipe.Set(ctx, key, data, ttl)
	}
	
	_, err := pipe.Exec(ctx)
	return err
}