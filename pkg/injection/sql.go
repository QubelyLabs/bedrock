package injection

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	gormTxContextKey = "gorm_tx_context"
)

func SetSQL(c *gin.Context, v *gorm.DB) {
	c.Set(gormTxContextKey, v)
}

func GetSQL(c *gin.Context) *gorm.DB {
	tx := c.MustGet(gormTxContextKey)

	v := tx.(*gorm.DB)
	return v
}
