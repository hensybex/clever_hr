// internal/usecase/vacancy.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/utils"
	"github.com/pgvector/pgvector-go"
)

type VacancyUsecase interface {
	CreateVacancy(vacancy *model.Vacancy) error
	GetVacancyByID(id uint) (*model.Vacancy, error) // New method
	GetAllVacancies() ([]model.Vacancy, error)      // New method
	UpdateVacancyStatus(id uint, status string) error
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

func (uc *vacancyUsecaseImpl) CreateVacancy(vacancy *model.Vacancy) error {
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
	vacancyEmbedding := &model.VacancyEmbedding{
		VacancyID: vacancy.ID,
		Embedding: pgvector.NewVector(embeddingFloat32), // Correct conversion
	}

	if err := uc.embeddingRepo.CreateVacancyEmbedding(vacancyEmbedding); err != nil {
		return err
	}

	return nil
}

func (uc *vacancyUsecaseImpl) GetVacancyByID(id uint) (*model.Vacancy, error) {
	return uc.vacancyRepo.GetByID(id)
}

func (uc *vacancyUsecaseImpl) GetAllVacancies() ([]model.Vacancy, error) {
	return uc.vacancyRepo.GetAll()
}

func (uc *vacancyUsecaseImpl) UpdateVacancyStatus(id uint, status string) error {
	return uc.vacancyRepo.UpdateStatus(id, status)
}
