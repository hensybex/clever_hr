package repository

import "clever_hr_embeddings/internal/models"
import "gorm.io/gorm"

type VacancyRepository interface {
	Create(vacancy *models.Vacancy) error
	GetByID(id uint) (*models.Vacancy, error)
}

type vacancyRepositoryImpl struct {
	db *gorm.DB
}

func NewVacancyRepository(db *gorm.DB) VacancyRepository {
	return &vacancyRepositoryImpl{db}
}

func (repo *vacancyRepositoryImpl) Create(vacancy *models.Vacancy) error {
	return repo.db.Create(vacancy).Error
}

func (repo *vacancyRepositoryImpl) GetByID(id uint) (*models.Vacancy, error) {
	var vacancy models.Vacancy
	if err := repo.db.First(&vacancy, id).Error; err != nil {
		return nil, err
	}
	return &vacancy, nil
}
