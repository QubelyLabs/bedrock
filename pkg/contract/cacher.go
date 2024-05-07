package contract

import "time"

type Cacher interface {
	Get(key string) (any, error)
	Set(key string, value any, exp time.Duration) error
	Cache(key string, fn func() (any, error), exp time.Duration) (any, error)
}
