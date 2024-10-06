// internal/repository/interview_repository.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type InterviewRepository interface {
	CreateInterview(interview *model.Interview) error
	GetInterviewByID(id uint) (*model.Interview, error)
	UpdateInterviewStatus(id uint, status string) error
	GetInterviewByResumeID(resumeID uint) (*model.Interview, error)
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

func (r *interviewRepository) UpdateInterviewStatus(id uint, status string) error {
	return r.db.Model(&model.Interview{}).Where("id = ?", id).Update("status", status).Error
}

func (r *interviewRepository) GetInterviewByResumeID(resumeID uint) (*model.Interview, error) {
	var interview model.Interview
	err := r.db.Where("resume_id = ?", resumeID).Order("created_at DESC").First(&interview).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}
