// internal/repository/interview_repository.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type InterviewRepository interface {
	CreateInterview(interview *model.Interview) error
	GetInterviewByID(id uint) (*model.Interview, error)
}

type interviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db}
}

func (r *interviewRepository) CreateInterview(interview *model.Interview) error {
	return r.db.Create(interview).Error
}

func (r *interviewRepository) GetInterviewByID(id uint) (*model.Interview, error) {
	var interview model.Interview
	err := r.db.First(&interview, id).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}
