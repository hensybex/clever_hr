package repository

import "clever_hr_embeddings/internal/models"
import "gorm.io/gorm"

type ResumeRepository interface {
	Create(resume *models.Resume) error
	GetByID(id uint) (*models.Resume, error)
}

type resumeRepositoryImpl struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) ResumeRepository {
	return &resumeRepositoryImpl{db}
}

func (repo *resumeRepositoryImpl) Create(resume *models.Resume) error {
	return repo.db.Create(resume).Error
}

func (repo *resumeRepositoryImpl) GetByID(id uint) (*models.Resume, error) {
	var resume models.Resume
	if err := repo.db.First(&resume, id).Error; err != nil {
		return nil, err
	}
	return &resume, nil
}
