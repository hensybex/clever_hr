// repository/job_group.go

package repository

import (
	"clever_hr_api/internal/model/categories"
	"gorm.io/gorm"
)

type JobGroupRepository interface {
	Create(group *category_model.JobGroup) error
	GetByID(id int) (*category_model.JobGroup, error)
	GetAll() ([]category_model.JobGroup, error)
	Update(id int, updatedGroup *category_model.JobGroup) error
}

type jobGroupRepositoryImpl struct {
	db *gorm.DB
}

func NewJobGroupRepository(db *gorm.DB) JobGroupRepository {
	return &jobGroupRepositoryImpl{db}
}

func (repo *jobGroupRepositoryImpl) Create(group *category_model.JobGroup) error {
	return repo.db.Create(group).Error
}

func (repo *jobGroupRepositoryImpl) GetByID(id int) (*category_model.JobGroup, error) {
	var group category_model.JobGroup
	if err := repo.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (repo *jobGroupRepositoryImpl) GetAll() ([]category_model.JobGroup, error) {
	var groups []category_model.JobGroup
	if err := repo.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (repo *jobGroupRepositoryImpl) Update(id int, updatedGroup *category_model.JobGroup) error {
	return repo.db.Model(&category_model.JobGroup{}).Where("id = ?", id).Updates(updatedGroup).Error
}
