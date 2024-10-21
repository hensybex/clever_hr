// internal/model/resume_embedding.go

package model

import (
	"gorm.io/gorm"
	"time"
)

// ResumeEmbedding stores the embedding vector for a resume.
type ResumeEmbedding struct {
	ID                uint      `gorm:"primaryKey"`
	ResumeID          uint      `gorm:"uniqueIndex"`
	Embedding         []float32 `gorm:"type:float[]"` // Updated to use float32 slice
	OriginalEmbedding []float32 `gorm:"type:float[]"` // Updated to use float32 slice
	Resume            Resume    `gorm:"foreignKey:ResumeID"`
	Score             float64   `gorm:"score"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
