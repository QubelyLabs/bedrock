package db

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var kv *redis.Client

func InitKV(host string, port int, password string, db int) {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: password,
		DB:       db,
	})

	kv = r
}

func KV() *redis.Client {
	return kv
}
