// repository/resume.go

package repository

import (
	"clever_hr_api/internal/model"

	"gorm.io/gorm"
)

type ResumeRepository interface {
	CreateResume(resume *model.Resume) error
	GetByID(id uint) (*model.Resume, error)
	GetResumeByCandidateID(candidateID uint) ([]*model.Resume, error)
	GetAll() ([]model.Resume, error)
}

type resumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) ResumeRepository {
	return &resumeRepository{db}
}

func (r *resumeRepository) CreateResume(resume *model.Resume) error {
	return r.db.Create(resume).Error
}

func (r *resumeRepository) GetByID(id uint) (*model.Resume, error) {
	var resume model.Resume
	err := r.db.First(&resume, id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *resumeRepository) GetResumeByCandidateID(candidateID uint) ([]*model.Resume, error) {
	var resumes []*model.Resume
	err := r.db.Where("candidate_id = ?", candidateID).Find(&resumes).Error
	if err != nil {
		return nil, err
	}
	return resumes, nil
}

func (r *resumeRepository) GetAll() ([]model.Resume, error) {
	var resumes []model.Resume
	if err := r.db.Find(&resumes).Error; err != nil {
		return nil, err
	}
	return resumes, nil
}
