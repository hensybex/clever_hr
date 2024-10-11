package models

// Resume represents a candidate's resume.
type Resume struct {
	ID             uint `gorm:"primaryKey"`
	Name           string
	Email          string
	PhoneNumber    string
	Content        string
	Skills         []string `gorm:"type:text[]"`
	Experience     []string `gorm:"type:text[]"`
	Certifications []string `gorm:"type:text[]"`
	Education      []string `gorm:"type:text[]"`
}
