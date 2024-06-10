// Package cache provides an abstraction for caching operations using various cache implementations.
// currently, only a redis bases cache is implemented
package cache

import (
	"time"

	"github.com/QubelyLabs/bedrock/pkg/contract"
)

func NewDefaultCache() contract.Cacher {
	return NewRedisCache()
}

func Cache(key string, fn func() (any, error), exp time.Duration) (any, error) {
	cache := NewDefaultCache()
	return cache.Cache(key, fn, exp)
}
