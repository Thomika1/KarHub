package crud

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// BaseModel is a base model for all models
type BaseModel struct {
	CreatedAt *time.Time      `json:"createdAt,omitempty" example:"2024-06-01T00:00:00Z"`
	UpdatedAt *time.Time      `json:"updatedAt,omitempty" example:"2024-06-03T00:00:00Z"`
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"index" example:"2024-06-05T00:00:00Z"`
	ID        *string         `json:"id,omitempty" gorm:"primarykey;type:varchar(20)" example:"0613b3b4-4b3b-4b3b-4b3b-4b3b4b3b4b3b"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == nil || *b.ID == "" {
		id := xid.New().String()
		b.ID = &id
	}

	return nil
}
