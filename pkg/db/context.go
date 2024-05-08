package db

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	gormTxContextKey = "gorm_tx_context"
)

func SetSQLToContext(ctx *gin.Context, tx *gorm.DB) {
	ctx.Set(gormTxContextKey, tx)
}

func GetSQLFromContext(ctx *gin.Context) *gorm.DB {
	tx := ctx.MustGet(gormTxContextKey)

	gormTx := tx.(*gorm.DB)
	return gormTx
}
