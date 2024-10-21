// internal/model/vacancy.go

package model

import category_model "clever_hr_api/internal/model/categories"

// Vacancy represents a job vacancy.
type Vacancy struct {
	ID               uint `gorm:"primaryKey"`
	UploaderID       uint `json:"user_id"`
	Title            string
	Description      string
	StandarizedText  string
	JobGroupID       *uint `gorm:"index"`
	SpecializationID *uint `gorm:"index"`
	QualificationID  *uint `gorm:"index"`
	GPTCallID        *uint
	JobGroup         *category_model.JobGroup       `gorm:"foreignKey:JobGroupID"`
	Specialization   *category_model.Specialization `gorm:"foreignKey:SpecializationID"`
	Qualification    *category_model.Qualification  `gorm:"foreignKey:QualificationID"`
	GPTCall          *GPTCall                       `gorm:"foreignKey:GPTCallID"`
}
