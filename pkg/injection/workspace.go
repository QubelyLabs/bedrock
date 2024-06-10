package injection

import (
	"github.com/QubelyLabs/bedrock/pkg/util"
	"github.com/gin-gonic/gin"
)

const (
	workspaceContextKey = "workspace_context"
)

func SetWorkspace(c *gin.Context, v util.Object) {
	c.Set(workspaceContextKey, v)
}

func GetWorkspace(c *gin.Context) util.Object {
	tx := c.MustGet(workspaceContextKey)

	v := tx.(util.Object)
	return v
}
