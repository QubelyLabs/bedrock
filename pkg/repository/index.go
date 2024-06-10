package repository

import (
	"github.com/QubelyLabs/bedrock/pkg/injection"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository[E any] struct{}

func (r *Repository[E]) SQL(c *gin.Context) *gorm.DB {
	return injection.GetSQL(c)
}

func (r *Repository[E]) UpsertOne(c *gin.Context, entity *E) error {
	err := r.SQL(c).WithContext(c.Request.Context()).Save(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) UpsertMany(c *gin.Context, entities ...E) error {
	err := r.SQL(c).WithContext(c.Request.Context()).Save(entities).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) CreateOne(c *gin.Context, entity *E) error {
	err := r.SQL(c).WithContext(c.Request.Context()).Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) CreateMany(c *gin.Context, entities ...E) error {
	err := r.SQL(c).WithContext(c.Request.Context()).Create(entities).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) UpdateOne(c *gin.Context, id string, entity *E) error {
	err := r.SQL(c).WithContext(c.Request.Context()).Where("id = ?", id).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) UpdateMany(c *gin.Context, entity *E, query any, args ...any) error {
	err := r.SQL(c).WithContext(c.Request.Context()).Where(query, args...).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) FindOne(c *gin.Context, id string) (E, error) {
	entity := new(E)
	err := r.SQL(c).WithContext(c.Request.Context()).Where("id = ?", id).First(entity).Error
	if err != nil {
		return *entity, err
	}

	return *entity, nil
}

func (r *Repository[E]) FindMany(c *gin.Context, query any, args ...any) ([]E, error) {
	return r.FindManyWithLimit(c, -1, -1, query, args...)
}

func (r *Repository[E]) FindAll(c *gin.Context) ([]E, error) {
	return r.FindManyWithLimit(c, -1, -1, nil)
}

func (r *Repository[E]) FindManyWithLimit(c *gin.Context, limit int, offset int, query any, args ...any) ([]E, error) {
	entities := new([]E)
	err := r.SQL(c).WithContext(c.Request.Context()).Where(query, args...).Limit(limit).Offset(offset).Find(entities).Error
	if err != nil {
		return nil, err
	}

	return *entities, nil
}

func (r *Repository[E]) DeleteOne(c *gin.Context, id string) error {
	entity := new(E)
	err := r.SQL(c).WithContext(c.Request.Context()).Where("id = ?", id).Delete(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) DeleteMany(c *gin.Context, query any, args ...any) error {
	entity := new(E)
	err := r.SQL(c).WithContext(c.Request.Context()).Where(query, args...).Delete(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) Count(c *gin.Context, query any, args ...any) (i int64, err error) {
	entity := new(E)
	err = r.SQL(c).WithContext(c.Request.Context()).Where(query, args...).Model(entity).Count(&i).Error
	return
}

func NewRepository[E any]() *Repository[E] {
	return &Repository[E]{}
}
