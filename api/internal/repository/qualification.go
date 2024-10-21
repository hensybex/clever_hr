// repository/qualification.go

package repository

import (
	"clever_hr_api/internal/model/categories"
	"gorm.io/gorm"
)

type QualificationRepository interface {
	Create(qualification *category_model.Qualification) error
	GetByID(id int) (*category_model.Qualification, error)
	GetAll() ([]category_model.Qualification, error)
	Update(id int, updatedQualification *category_model.Qualification) error
}

type qualificationRepositoryImpl struct {
	db *gorm.DB
}

func NewQualificationRepository(db *gorm.DB) QualificationRepository {
	return &qualificationRepositoryImpl{db}
}

func (repo *qualificationRepositoryImpl) Create(qualification *category_model.Qualification) error {
	return repo.db.Create(qualification).Error
}

func (repo *qualificationRepositoryImpl) GetByID(id int) (*category_model.Qualification, error) {
	var qualification category_model.Qualification
	if err := repo.db.First(&qualification, id).Error; err != nil {
		return nil, err
	}
	return &qualification, nil
}

func (repo *qualificationRepositoryImpl) GetAll() ([]category_model.Qualification, error) {
	var qualifications []category_model.Qualification
	if err := repo.db.Find(&qualifications).Error; err != nil {
		return nil, err
	}
	return qualifications, nil
}

func (repo *qualificationRepositoryImpl) Update(id int, updatedQualification *category_model.Qualification) error {
	return repo.db.Model(&category_model.Qualification{}).Where("id = ?", id).Updates(updatedQualification).Error
}
