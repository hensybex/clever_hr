// repository/resume_analysis_result.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type ResumeAnalysisResultRepository interface {
	CreateResumeAnalysisResult(result *model.ResumeAnalysisResult) error
	GetResumeAnalysisResultByResumeID(resumeID uint) (*model.ResumeAnalysisResult, error)
	WasResumeAnalysed(resumeID uint) (bool, error)
}

type resumeAnalysisResultRepository struct {
	db *gorm.DB
}

func NewResumeAnalysisResultRepository(db *gorm.DB) ResumeAnalysisResultRepository {
	return &resumeAnalysisResultRepository{db}
}

func (r *resumeAnalysisResultRepository) CreateResumeAnalysisResult(result *model.ResumeAnalysisResult) error {
	return r.db.Create(result).Error
}

func (r *resumeAnalysisResultRepository) GetResumeAnalysisResultByResumeID(resumeID uint) (*model.ResumeAnalysisResult, error) {
	var result model.ResumeAnalysisResult
	err := r.db.Where("resume_id = ?", resumeID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *resumeAnalysisResultRepository) WasResumeAnalysed(resumeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ResumeAnalysisResult{}).Where("resume_id = ?", resumeID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
