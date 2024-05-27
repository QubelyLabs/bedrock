package injection

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	redisContextKey = "redis_context"
)

func SetKV(c *gin.Context, v *redis.Client) {
	c.Set(redisContextKey, v)
}

func GetKV(c *gin.Context) *redis.Client {
	tx := c.MustGet(redisContextKey)

	v := tx.(*redis.Client)
	return v
}
