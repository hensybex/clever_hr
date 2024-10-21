// model/gpt_call.go

package model

import (
	"gorm.io/gorm"
	"time"
)

type GPTCall struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id,omitempty" swaggerignore:"true"`
	CreatedAt   time.Time      `json:"createdAt,omitempty" swaggerignore:"true"`
	UpdatedAt   time.Time      `json:"updatedAt,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" swaggerignore:"true"`
	FinalPrompt string         `json:"finalPrompt,omitempty"`
	Reply       string         `json:"reply,omitempty"`
	EntityType  string         `json:"entity_type,omitempty"`
	EntityID    uint           `json:"entity_id,omitempty"`
}
