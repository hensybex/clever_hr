package usecases

import (
	"clever_hr_embeddings/internal/models"
	"clever_hr_embeddings/internal/repository"
	"gorm.io/gorm"
)

type CandidateMatchingUsecase interface {
	GetBestCandidates(vacancyID uint) ([]models.Resume, error)
}

type candidateMatchingUsecaseImpl struct {
	db            *gorm.DB
	resumeRepo    repository.ResumeRepository
	vacancyRepo   repository.VacancyRepository
	embeddingRepo repository.EmbeddingRepository
}

func NewCandidateMatchingUsecase(
	resumeRepo repository.ResumeRepository,
	vacancyRepo repository.VacancyRepository,
	embeddingRepo repository.EmbeddingRepository,
) CandidateMatchingUsecase {
	return &candidateMatchingUsecaseImpl{
		resumeRepo:    resumeRepo,
		vacancyRepo:   vacancyRepo,
		embeddingRepo: embeddingRepo,
	}
}

func (uc *candidateMatchingUsecaseImpl) GetBestCandidates(vacancyID uint) ([]models.Resume, error) {
	// Retrieve the embedding for the given vacancy
	var vacancyEmbedding models.VacancyEmbedding
	if err := uc.embeddingRepo.GetVacancyEmbeddingByVacancyID(vacancyID, &vacancyEmbedding); err != nil {
		return nil, err
	}

	// Find resumes with similar embeddings using pgvector
	var resumes []models.Resume
	err := uc.db.Raw(`
        SELECT r.*
        FROM resume_embeddings re
        JOIN resumes r ON re.resume_id = r.id
        ORDER BY re.embedding <-> ? -- Cosine similarity or Euclidean distance
        LIMIT 10`, vacancyEmbedding.Embedding).
		Scan(&resumes).Error
	if err != nil {
		return nil, err
	}

	return resumes, nil
}
