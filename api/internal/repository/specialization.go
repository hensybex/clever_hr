// repository/specialization.go

package repository

import (
	"clever_hr_api/internal/model/categories"
	"gorm.io/gorm"
)

type SpecializationRepository interface {
	Create(specialization *category_model.Specialization) error
	GetByID(id int) (*category_model.Specialization, error)
	GetAll() ([]category_model.Specialization, error)
	GetAllByGroupID(groupID int) ([]category_model.Specialization, error)
	Update(id int, updatedSpecialization *category_model.Specialization) error
}

type specializationRepositoryImpl struct {
	db *gorm.DB
}

func NewSpecializationRepository(db *gorm.DB) SpecializationRepository {
	return &specializationRepositoryImpl{db}
}

func (repo *specializationRepositoryImpl) Create(specialization *category_model.Specialization) error {
	return repo.db.Create(specialization).Error
}

func (repo *specializationRepositoryImpl) GetByID(id int) (*category_model.Specialization, error) {
	var specialization category_model.Specialization
	if err := repo.db.Preload("JobGroup").First(&specialization, id).Error; err != nil {
		return nil, err
	}
	return &specialization, nil
}

func (repo *specializationRepositoryImpl) GetAll() ([]category_model.Specialization, error) {
	var specializations []category_model.Specialization
	if err := repo.db.Preload("JobGroup").Find(&specializations).Error; err != nil {
		return nil, err
	}
	return specializations, nil
}

func (repo *specializationRepositoryImpl) GetAllByGroupID(groupID int) ([]category_model.Specialization, error) {
	var specializations []category_model.Specialization
	if err := repo.db.Where("job_group_id = ?", groupID).Preload("JobGroup").Find(&specializations).Error; err != nil {
		return nil, err
	}
	return specializations, nil
}

func (repo *specializationRepositoryImpl) Update(id int, updatedSpecialization *category_model.Specialization) error {
	return repo.db.Model(&category_model.Specialization{}).Where("id = ?", id).Updates(updatedSpecialization).Error
}
