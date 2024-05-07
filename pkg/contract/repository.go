package contract

import (
	"context"

	"gorm.io/gorm"
)

type Repository[E any] interface {
	DB() *gorm.DB
	UpsertOne(context.Context, *E) error
	UpsertMany(context.Context, ...E) error
	CreateOne(context.Context, *E) error
	CreateMany(context.Context, ...E) error
	UpdateOne(context.Context, string, *E) error
	UpdateMany(context.Context, *E, any, ...any) error
	FindOne(context.Context, string) (E, error)
	FindMany(context.Context, any, ...any) ([]E, error)
	FindAll(context.Context) ([]E, error)
	FindManyWithLimit(context.Context, int, int, any, ...any) ([]E, error)
	DeleteOne(context.Context, string) error
	DeleteMany(context.Context, any, ...any) error
	Count(context.Context, any, ...any) (int64, error)
}
