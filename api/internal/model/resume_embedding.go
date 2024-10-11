// internal/model/resume_embedding.go

package model

import "github.com/pgvector/pgvector-go"

// ResumeEmbedding stores the embedding vector for a resume.
type ResumeEmbedding struct {
	ID        uint            `gorm:"primaryKey"`
	ResumeID  uint            `gorm:"uniqueIndex"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)"`
	Resume    Resume          `gorm:"foreignKey:ResumeID"`
}
