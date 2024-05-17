package controller

import (
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

func (ctrl *BaseController) Success(c *gin.Context, message string, data any) {
	c.JSON(200, gin.H{
		"status":  true,
		"message": message,
		"data":    data,
	})
}

func (ctrl *BaseController) Error(c *gin.Context, message string) {
	c.JSON(400, gin.H{
		"status":  false,
		"message": message,
	})
}

func (ctrl *BaseController) ErrorWithData(c *gin.Context, message string, data any) {
	c.JSON(400, gin.H{
		"status":  false,
		"message": message,
		"data":    data,
	})
}

func (ctrl *BaseController) ErrorWithCode(c *gin.Context, message string, code int) {
	c.JSON(400, gin.H{
		"status":  false,
		"message": message,
	})
}

func (ctrl *BaseController) ErrorWithDataAndCode(c *gin.Context, message string, data any, code int) {
	if code == 0 {
		code = 400
	}

	c.JSON(code, gin.H{
		"status":  false,
		"message": message,
		"data":    data,
	})
}
