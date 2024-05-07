package repository

import (
	"context"

	"github.com/qubelylabs/bedrock/pkg/db"
	"gorm.io/gorm"
)

type Repository[E any] struct {
	db *gorm.DB
}

func (r *Repository[E]) DB() *gorm.DB {
	return r.db
}

func (r *Repository[E]) UpsertOne(ctx context.Context, entity *E) error {
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) UpsertMany(ctx context.Context, entities ...E) error {
	err := r.db.WithContext(ctx).Save(entities).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) CreateOne(ctx context.Context, entity *E) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) CreateMany(ctx context.Context, entities ...E) error {
	err := r.db.WithContext(ctx).Create(entities).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) UpdateOne(ctx context.Context, id string, entity *E) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) UpdateMany(ctx context.Context, entity *E, query any, args ...any) error {
	err := r.db.WithContext(ctx).Where(query, args...).Updates(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) FindOne(ctx context.Context, id string) (E, error) {
	entity := new(E)
	err := r.db.WithContext(ctx).Where("id = ?", id).First(entity).Error
	if err != nil {
		return *entity, err
	}

	return *entity, nil
}

func (r *Repository[E]) FindMany(ctx context.Context, query any, args ...any) ([]E, error) {
	return r.FindManyWithLimit(ctx, -1, -1, query, args...)
}

func (r *Repository[E]) FindAll(ctx context.Context) ([]E, error) {
	return r.FindManyWithLimit(ctx, -1, -1, nil)
}

func (r *Repository[E]) FindManyWithLimit(ctx context.Context, limit int, offset int, query any, args ...any) ([]E, error) {
	entities := new([]E)
	err := r.db.WithContext(ctx).Where(query, args...).Limit(limit).Offset(offset).Find(entities).Error
	if err != nil {
		return nil, err
	}

	return *entities, nil
}

func (r *Repository[E]) DeleteOne(ctx context.Context, id string) error {
	entity := new(E)
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) DeleteMany(ctx context.Context, query any, args ...any) error {
	entity := new(E)
	err := r.db.WithContext(ctx).Where(query, args...).Delete(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[E]) Count(ctx context.Context, query any, args ...any) (i int64, err error) {
	entity := new(E)
	err = r.db.WithContext(ctx).Where(query, args...).Model(entity).Count(&i).Error
	return
}

func NewRepository[E any]() *Repository[E] {
	return &Repository[E]{
		db.SQL(),
	}
}
