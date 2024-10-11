package repository

import "clever_hr_embeddings/internal/models"
import "gorm.io/gorm"

type EmbeddingRepository interface {
	CreateResumeEmbedding(embedding *models.ResumeEmbedding) error
	CreateVacancyEmbedding(embedding *models.VacancyEmbedding) error
	GetVacancyEmbeddingByVacancyID(vacancyID uint, embedding *models.VacancyEmbedding) error
	GetAllResumeEmbeddings() ([]models.ResumeEmbedding, error)
}

type embeddingRepositoryImpl struct {
	db *gorm.DB
}

func NewEmbeddingRepository(db *gorm.DB) EmbeddingRepository {
	return &embeddingRepositoryImpl{db}
}

func (repo *embeddingRepositoryImpl) CreateResumeEmbedding(embedding *models.ResumeEmbedding) error {
	return repo.db.Create(embedding).Error
}

func (repo *embeddingRepositoryImpl) CreateVacancyEmbedding(embedding *models.VacancyEmbedding) error {
	return repo.db.Create(embedding).Error
}

func (repo *embeddingRepositoryImpl) GetVacancyEmbeddingByVacancyID(vacancyID uint, embedding *models.VacancyEmbedding) error {
	return repo.db.First(&embedding, "vacancy_id = ?", vacancyID).Error
}

func (repo *embeddingRepositoryImpl) GetAllResumeEmbeddings() ([]models.ResumeEmbedding, error) {
	var embeddings []models.ResumeEmbedding
	err := repo.db.Find(&embeddings).Error
	return embeddings, err
}
