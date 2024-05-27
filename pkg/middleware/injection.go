package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qubelylabs/bedrock/pkg/injection"
	"github.com/qubelylabs/bedrock/pkg/util"
)

func Injection() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("x-user")
		if user != "" {
			injection.SetUser(c, util.FromBase64(user))
		}

		workspace := c.GetHeader("x-workspace")
		if workspace != "" {
			injection.SetWorkspace(c, util.FromBase64(workspace))
		}

		c.Next()
	}
}
