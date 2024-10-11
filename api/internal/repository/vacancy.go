// internal/repository/vacancy.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type VacancyRepository interface {
	Create(vacancy *model.Vacancy) error
	GetByID(id uint) (*model.Vacancy, error)
	GetAll() ([]model.Vacancy, error)          // New method
	UpdateStatus(id uint, status string) error // New method
}

type vacancyRepositoryImpl struct {
	db *gorm.DB
}

func NewVacancyRepository(db *gorm.DB) VacancyRepository {
	return &vacancyRepositoryImpl{db}
}

func (repo *vacancyRepositoryImpl) Create(vacancy *model.Vacancy) error {
	return repo.db.Create(vacancy).Error
}

func (repo *vacancyRepositoryImpl) GetByID(id uint) (*model.Vacancy, error) {
	var vacancy model.Vacancy
	if err := repo.db.First(&vacancy, id).Error; err != nil {
		return nil, err
	}
	return &vacancy, nil
}

func (repo *vacancyRepositoryImpl) GetAll() ([]model.Vacancy, error) {
	var vacancies []model.Vacancy
	if err := repo.db.Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}

func (repo *vacancyRepositoryImpl) UpdateStatus(id uint, status string) error {
	return repo.db.Model(&model.Vacancy{}).Where("id = ?", id).Update("status", status).Error
}
