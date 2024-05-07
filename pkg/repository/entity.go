package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	ID        string         `gorm:"primaryKey;column:id;type:string;size:36"`
	CreatedAt time.Time      `gorm:"index;column:created_at"`
	UpdatedAt time.Time      `gorm:"index;column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New().String()
	return nil
}
