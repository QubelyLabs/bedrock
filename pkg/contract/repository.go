package contract

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository[E any] interface {
	SQL(*gin.Context) *gorm.DB
	UpsertOne(*gin.Context, *E) error
	UpsertMany(*gin.Context, ...E) error
	CreateOne(*gin.Context, *E) error
	CreateMany(*gin.Context, ...E) error
	UpdateOne(*gin.Context, string, *E) error
	UpdateMany(*gin.Context, *E, any, ...any) error
	FindOne(*gin.Context, string) (E, error)
	FindMany(*gin.Context, any, ...any) ([]E, error)
	FindAll(*gin.Context) ([]E, error)
	FindManyWithLimit(*gin.Context, int, int, any, ...any) ([]E, error)
	DeleteOne(*gin.Context, string) error
	DeleteMany(*gin.Context, any, ...any) error
	Count(*gin.Context, any, ...any) (int64, error)
}
