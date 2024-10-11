// internal/repository/embedding.go

package repository

import "clever_hr_api/internal/model"
import "gorm.io/gorm"

type EmbeddingRepository interface {
	CreateResumeEmbedding(embedding *model.ResumeEmbedding) error
	CreateVacancyEmbedding(embedding *model.VacancyEmbedding) error
	GetVacancyEmbeddingByVacancyID(vacancyID uint, embedding *model.VacancyEmbedding) error
	GetAllResumeEmbeddings() ([]model.ResumeEmbedding, error)
}

type embeddingRepositoryImpl struct {
	db *gorm.DB
}

func NewEmbeddingRepository(db *gorm.DB) EmbeddingRepository {
	return &embeddingRepositoryImpl{db}
}

func (repo *embeddingRepositoryImpl) CreateResumeEmbedding(embedding *model.ResumeEmbedding) error {
	return repo.db.Create(embedding).Error
}

func (repo *embeddingRepositoryImpl) CreateVacancyEmbedding(embedding *model.VacancyEmbedding) error {
	return repo.db.Create(embedding).Error
}

func (repo *embeddingRepositoryImpl) GetVacancyEmbeddingByVacancyID(vacancyID uint, embedding *model.VacancyEmbedding) error {
	return repo.db.First(&embedding, "vacancy_id = ?", vacancyID).Error
}

func (repo *embeddingRepositoryImpl) GetAllResumeEmbeddings() ([]model.ResumeEmbedding, error) {
	var embeddings []model.ResumeEmbedding
	err := repo.db.Find(&embeddings).Error
	return embeddings, err
}
