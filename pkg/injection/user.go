package injection

import (
	"github.com/QubelyLabs/bedrock/pkg/util"
	"github.com/gin-gonic/gin"
)

const (
	userContextKey = "user_context"
)

func SetUser(c *gin.Context, v util.Object) {
	c.Set(userContextKey, v)
}

func GetUser(c *gin.Context) util.Object {
	tx := c.MustGet(userContextKey)

	v := tx.(util.Object)
	return v
}
