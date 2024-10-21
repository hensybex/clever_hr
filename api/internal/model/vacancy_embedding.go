// internal/model/vacancy_embedding.go

package model

import (
	"gorm.io/gorm"
	"time"
)

// VacancyEmbedding stores the embedding vector for a vacancy.
type VacancyEmbedding struct {
	ID                uint      `gorm:"primaryKey"`
	VacancyID         uint      `gorm:"uniqueIndex"`
	Embedding         []float32 `gorm:"type:float[]"` // Updated to use float32 slice
	OriginalEmbedding []float32 `gorm:"type:float[]"` // Updated to use float32 slice
	Vacancy           Vacancy   `gorm:"foreignKey:VacancyID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
