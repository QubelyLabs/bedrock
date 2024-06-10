package controller

import (
	"errors"
	"log"

	"github.com/QubelyLabs/bedrock/pkg/injection"
	"github.com/QubelyLabs/bedrock/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BaseController struct{}

func (ctrl *BaseController) SQL(c *gin.Context) *gorm.DB {
	return injection.GetSQL(c)
}

func (ctrl *BaseController) KV(c *gin.Context) *redis.Client {
	return injection.GetKV(c)
}

func (ctrl *BaseController) User(c *gin.Context) (util.Object, error) {
	user := injection.GetUser(c)
	userId, ok := user["id"].(string)
	if !ok || userId == "" {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (ctrl *BaseController) Workspace(c *gin.Context) (util.Object, error) {
	workspace := injection.GetWorkspace(c)
	workspaceId, ok := workspace["id"].(string)
	if !ok || workspaceId == "" {
		return nil, errors.New("workspace not found")
	}

	return workspace, nil
}

func (ctrl *BaseController) Validate(c *gin.Context, payload any) (any, bool) {
	err := c.ShouldBindJSON(payload)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
			return nil, false
		}

		if data, ok := err.(validator.ValidationErrors); ok {
			log.Println(err)
			return data, false
		}

		log.Println(err)
		return nil, false
	}

	return nil, true
}

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
