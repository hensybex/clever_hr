package models

// Vacancy represents a job vacancy.
type Vacancy struct {
	ID                uint `gorm:"primaryKey"`
	Title             string
	Description       string
	PriceRange        string
	YearsOfExperience int
	WorkType          string
	Requirements      []string `gorm:"type:text[]"`
	Qualifications    []string `gorm:"type:text[]"`
}
