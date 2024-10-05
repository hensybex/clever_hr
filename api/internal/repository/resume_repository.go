// internal/repository/resume_repository.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type ResumeRepository interface {
	CreateResume(resume *model.Resume) error
	GetResumeByID(id uint) (*model.Resume, error)
	GetResumeByTgID(tgID uint) (*model.Resume, error)
	GetResumeByCandidateID(candidateID uint) (*model.Resume, error)
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

func (r *resumeRepository) GetResumeByID(id uint) (*model.Resume, error) {
	var resume model.Resume
	err := r.db.First(&resume, id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *resumeRepository) GetResumeByTgID(tgID uint) (*model.Resume, error) {
	var resume model.Resume
	err := r.db.Where("candidate_id = ?", tgID).First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *resumeRepository) GetResumeByCandidateID(candidateID uint) (*model.Resume, error) {
	var resume model.Resume
	err := r.db.Where("candidate_id = ?", candidateID).First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}
