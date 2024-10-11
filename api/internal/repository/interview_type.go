// internal/repository/interview_type.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type InterviewTypeRepository interface {
	GetInterviewTypeByID(id uint) (*model.InterviewType, error)
	ListInterviewTypes() ([]model.InterviewType, error)
}

type interviewTypeRepository struct {
	db *gorm.DB
}

func NewInterviewTypeRepository(db *gorm.DB) InterviewTypeRepository {
	return &interviewTypeRepository{db}
}

func (r *interviewTypeRepository) GetInterviewTypeByID(id uint) (*model.InterviewType, error) {
	var interviewType model.InterviewType
	err := r.db.First(&interviewType, id).Error
	if err != nil {
		return nil, err
	}
	return &interviewType, nil
}

func (r *interviewTypeRepository) ListInterviewTypes() ([]model.InterviewType, error) {
	var interviewTypes []model.InterviewType
	if err := r.db.Find(&interviewTypes).Error; err != nil {
		return nil, err
	}
	return interviewTypes, nil
}
