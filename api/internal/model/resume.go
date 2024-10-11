// internal/model/resume.go

package model

import (
	"time"
)

type Resume struct {
	ID              uint       `gorm:"primaryKey"`
	CandidateID     *uint      `gorm:"index"`
	Candidate       *Candidate `gorm:"foreignKey:CandidateID"`
	ResumePDF       string
	TextExtracted   string
	RewrittenResume string
	Name            *string
	Email           *string
	PhoneNumber     *string
	BirthDate       *time.Time
	TotalYears      *string
	PreferableJob   *string
	Skills          []string `gorm:"type:text[]"`
	Experience      []string `gorm:"type:text[]"`
	Certifications  []string `gorm:"type:text[]"`
	Education       []string `gorm:"type:text[]"`
	UploadedFrom    string
	GPTCallID       *uint    `gorm:"index"`
	GPTCall         *GPTCall `gorm:"foreignKey:GPTCallID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
