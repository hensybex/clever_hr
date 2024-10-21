// internal/usecase/embedding.go

package usecase

import (
	"clever_hr_api/internal/repository"
	"fmt"
	"log"
)

type EmbeddingUsecase interface {
	CreateResumeEmbedding(resumeID uint, embedding []float32) error
	CreateVacancyEmbedding(vacancyID uint, embedding []float32) error
	GetResumeEmbedding(resumeID uint) ([]float32, error)
	GetVacancyEmbedding(vacancyID uint) ([]float32, error)
	FindMatchingResumes(vacancyEmbedding []float32, limit int) ([]repository.ResumeMatch, error)
}

type embeddingUsecaseImpl struct {
	embeddingRepo repository.EmbeddingRepository
}

func NewEmbeddingUsecase(repo repository.EmbeddingRepository) EmbeddingUsecase {
	return &embeddingUsecaseImpl{repo}
}

func (uc *embeddingUsecaseImpl) CreateResumeEmbedding(resumeID uint, embedding []float32) error {
	return uc.embeddingRepo.CreateResumeEmbedding(resumeID, embedding)
}

func (uc *embeddingUsecaseImpl) CreateVacancyEmbedding(vacancyID uint, embedding []float32) error {
	return uc.embeddingRepo.CreateVacancyEmbedding(vacancyID, embedding)
}

func (uc *embeddingUsecaseImpl) GetResumeEmbedding(resumeID uint) ([]float32, error) {
	embedding, err := uc.embeddingRepo.GetResumeEmbedding(resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume embedding: %w", err)
	}
	return embedding, nil
}

func (uc *embeddingUsecaseImpl) GetVacancyEmbedding(vacancyID uint) ([]float32, error) {
	log.Printf("Starting GetVacancyEmbedding for vacancy ID: %d", vacancyID)

	embedding, err := uc.embeddingRepo.GetVacancyEmbedding(vacancyID)
	if err != nil {
		log.Printf("Error while getting vacancy embedding: %v", err)
		return nil, fmt.Errorf("failed to get vacancy embedding: %w", err)
	}

	log.Printf("Successfully retrieved embedding for vacancy ID %d: %v", vacancyID, embedding)
	return embedding, nil
}

func (uc *embeddingUsecaseImpl) FindMatchingResumes(vacancyEmbedding []float32, limit int) ([]repository.ResumeMatch, error) {
	return uc.embeddingRepo.FindMatchingResumes(vacancyEmbedding, limit)
}
