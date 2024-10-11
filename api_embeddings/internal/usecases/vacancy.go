package usecases

import (
	"clever_hr_embeddings/internal/models"
	"clever_hr_embeddings/internal/repository"
	"clever_hr_embeddings/internal/utils"
	"github.com/pgvector/pgvector-go"
)

type VacancyUsecase interface {
	CreateVacancy(vacancy *models.Vacancy) error
}

type vacancyUsecaseImpl struct {
	vacancyRepo   repository.VacancyRepository
	embeddingRepo repository.EmbeddingRepository
}

func NewVacancyUsecase(
	vacancyRepo repository.VacancyRepository,
	embeddingRepo repository.EmbeddingRepository,
) VacancyUsecase {
	return &vacancyUsecaseImpl{
		vacancyRepo:   vacancyRepo,
		embeddingRepo: embeddingRepo,
	}
}

func (uc *vacancyUsecaseImpl) CreateVacancy(vacancy *models.Vacancy) error {
	// Save vacancy to the database
	if err := uc.vacancyRepo.Create(vacancy); err != nil {
		return err
	}

	// Generate embedding for the vacancy description
	embedding, err := utils.GenerateEmbedding(vacancy.Description)
	if err != nil {
		return err
	}

	// Convert the []float64 embedding to []float32 for pgvector
	embeddingFloat32 := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingFloat32[i] = float32(v)
	}

	// Save the embedding as pgvector.Vector
	vacancyEmbedding := &models.VacancyEmbedding{
		VacancyID: vacancy.ID,
		Embedding: pgvector.NewVector(embeddingFloat32), // Correct conversion
	}

	if err := uc.embeddingRepo.CreateVacancyEmbedding(vacancyEmbedding); err != nil {
		return err
	}

	return nil
}
