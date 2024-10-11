package usecases

import (
	"clever_hr_embeddings/internal/models"
	"clever_hr_embeddings/internal/repository"
	"clever_hr_embeddings/internal/utils"

	"github.com/pgvector/pgvector-go"
)

type ResumeUsecase interface {
	CreateResume(resume *models.Resume) error
}

func NewResumeUsecase(
	resumeRepo repository.ResumeRepository,
	embeddingRepo repository.EmbeddingRepository,
) ResumeUsecase {
	return &resumeUsecaseImpl{
		resumeRepo:    resumeRepo,
		embeddingRepo: embeddingRepo,
	}
}

type resumeUsecaseImpl struct {
	resumeRepo    repository.ResumeRepository
	embeddingRepo repository.EmbeddingRepository
}

func (uc *resumeUsecaseImpl) CreateResume(resume *models.Resume) error {
	// Save resume to database
	if err := uc.resumeRepo.Create(resume); err != nil {
		return err
	}

	// Generate embedding
	embedding, err := utils.GenerateEmbedding(resume.Content)
	if err != nil {
		return err
	}

	// Convert the []float64 embedding to []float32 for pgvector
	embeddingFloat32 := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingFloat32[i] = float32(v)
	}

	// Save the embedding as pgvector.Vector
	resumeEmbedding := &models.ResumeEmbedding{
		ResumeID:  resume.ID,
		Embedding: pgvector.NewVector(embeddingFloat32), // Correct conversion
	}

	if err := uc.embeddingRepo.CreateResumeEmbedding(resumeEmbedding); err != nil {
		return err
	}

	return nil
}
