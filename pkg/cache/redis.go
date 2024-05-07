package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qubelylabs/bedrock/pkg/db"
	"github.com/redis/go-redis/v9"
)

// redisCache represents a Redis cache instance
type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

// Get retrieves data from the cache based on the provided key.
// It returns the retrieved data and an error, if any.
func (r *redisCache) Get(key string) (any, error) {
	str, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result any
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Set stores the provided value in the cache with the specified key with an expiration duration.
func (r *redisCache) Set(key string, value any, exp time.Duration) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = r.client.Set(r.ctx, key, string(buf), exp).Result()
	if err != nil {
		return err
	}

	return nil
}

// Cache provides a caching mechanism using the Redis cache.
// It retrieves data based on the provided key, and if it doesn't exist, it invokes the provided function to generate the data and stores it in the cache.
// It returns the retrieved or newly generated data and an error, if any.
func (r *redisCache) Cache(key string, fn func() (any, error), exp time.Duration) (any, error) {
	existing, err := r.Get(key)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == nil {
		return existing, nil
	}

	result, err := fn()
	if err != nil {
		return nil, err
	}

	err = r.Set(key, result, exp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewRedisCache() *redisCache {
	return &redisCache{db.KV(), context.Background()}
}
