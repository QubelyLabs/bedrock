package contract

import "github.com/gin-gonic/gin"

type Controller interface {
	CreateOne(c *gin.Context)
	CreateMany(c *gin.Context)
	UpdateOne(c *gin.Context)
	UpdateMany(c *gin.Context)
	FindOne(c *gin.Context)
	FindMany(c *gin.Context)
	DeleteOne(c *gin.Context)
	DeleteMany(c *gin.Context)
}
