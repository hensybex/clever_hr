// internal/repository/vacancy_resume_match.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type VacancyResumeMatchRepository interface {
	CreateMatch(vacancyResumeMatch *model.VacancyResumeMatch) error
	CreateOrUpdateMatch(vacancyResumeMatch *model.VacancyResumeMatch) error
	GetMatchesByVacancyID(vacancyID uint) ([]model.VacancyResumeMatch, error)
	GetMatchByID(matchID uint) (*model.VacancyResumeMatch, error)
}

type vacancyResumeMatchRepository struct {
	db *gorm.DB
}

func NewVacancyResumeMatchRepository(db *gorm.DB) VacancyResumeMatchRepository {
	return &vacancyResumeMatchRepository{db}
}

func (r *vacancyResumeMatchRepository) CreateMatch(vacancyResumeMatch *model.VacancyResumeMatch) error {
	return r.db.Create(vacancyResumeMatch).Error
}

func (r *vacancyResumeMatchRepository) CreateOrUpdateMatch(vacancyResumeMatch *model.VacancyResumeMatch) error {
	return r.db.
		Where("vacancy_id = ? AND resume_id = ?", vacancyResumeMatch.VacancyID, vacancyResumeMatch.ResumeID).
		Assign(vacancyResumeMatch).
		FirstOrCreate(vacancyResumeMatch).Error
}

// New method to get matches by vacancy ID
func (r *vacancyResumeMatchRepository) GetMatchesByVacancyID(vacancyID uint) ([]model.VacancyResumeMatch, error) {
	var matches []model.VacancyResumeMatch
	if err := r.db.Where("vacancy_id = ?", vacancyID).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *vacancyResumeMatchRepository) GetMatchByID(matchID uint) (*model.VacancyResumeMatch, error) {
	var match model.VacancyResumeMatch
	if err := r.db.First(&match, matchID).Error; err != nil {
		return nil, err
	}
	return &match, nil
}
