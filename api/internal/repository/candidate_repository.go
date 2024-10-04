// internal/repository/candidate_repository.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type CandidateRepository interface {
	CreateCandidate(candidate *model.Candidate) error
	GetCandidateByID(id uint) (*model.Candidate, error)
}

type candidateRepository struct {
	db *gorm.DB
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{db}
}

func (r *candidateRepository) CreateCandidate(candidate *model.Candidate) error {
	return r.db.Create(candidate).Error
}

func (r *candidateRepository) GetCandidateByID(id uint) (*model.Candidate, error) {
	var candidate model.Candidate
	err := r.db.First(&candidate, id).Error
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}
