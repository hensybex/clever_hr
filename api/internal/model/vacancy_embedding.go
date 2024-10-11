// internal/model/vacancy_embedding.go

package model

import "github.com/pgvector/pgvector-go"

// VacancyEmbedding stores the embedding vector for a vacancy.
type VacancyEmbedding struct {
	ID        uint            `gorm:"primaryKey"`
	VacancyID uint            `gorm:"uniqueIndex"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)"`
	Vacancy   Vacancy         `gorm:"foreignKey:VacancyID"`
}
