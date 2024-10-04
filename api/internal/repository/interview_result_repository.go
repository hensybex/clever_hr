// internal/repository/interview_result_repository.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type InterviewAnalysisResultRepository interface {
	CreateInterviewAnalysisResult(result *model.InterviewAnalysisResult) error
	GetInterviewAnalysisResultByInterviewID(interviewID uint) (*model.InterviewAnalysisResult, error)
}

type interviewResultRepository struct {
	db *gorm.DB
}

func NewInterviewAnalysisResultRepository(db *gorm.DB) InterviewAnalysisResultRepository {
	return &interviewResultRepository{db}
}

func (r *interviewResultRepository) CreateInterviewAnalysisResult(result *model.InterviewAnalysisResult) error {
	return r.db.Create(result).Error
}

func (r *interviewResultRepository) GetInterviewAnalysisResultByInterviewID(interviewID uint) (*model.InterviewAnalysisResult, error) {
	var result model.InterviewAnalysisResult
	err := r.db.Where("interview_id = ?", interviewID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
