package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	ID        string         `gorm:"primaryKey;column:id;type:string;size:36;not null" json:"id,omitempty"`
	CreatedAt *time.Time     `gorm:"index;column:created_at;not null" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"index;column:updated_at;not null;->:false;<-:create" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at;->:false;<-:create" json:"deleted_at,omitempty"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) error {
	if e.ID != "" {
		return nil
	}

	e.ID = uuid.New().String()
	now := time.Now()
	e.CreatedAt = &now
	e.UpdatedAt = &now
	return nil
}

func (e *Entity) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	e.UpdatedAt = &now
	return nil
}
