package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	ID        string         `gorm:"primaryKey;column:id;type:string;size:36;not null"`
	CreatedAt time.Time      `gorm:"index;column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"index;column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New().String()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Entity) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}
