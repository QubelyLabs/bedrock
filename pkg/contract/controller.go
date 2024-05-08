package contract

import "github.com/gin-gonic/gin"

type Controller interface {
	CreateOne(*gin.Context)
	CreateMany(*gin.Context)
	UpdateOne(*gin.Context)
	UpdateMany(*gin.Context)
	FindOne(*gin.Context)
	FindMany(*gin.Context)
	DeleteOne(*gin.Context)
	DeleteMany(*gin.Context)
}
