// internal/model/resume.go

package model

import (
	"time"
)

type Resume struct {
	ID               uint      `gorm:"primaryKey"`
	CandidateID      uint      `gorm:"index"`
	Candidate        Candidate `gorm:"foreignKey:CandidateID"`
	ResumePDF        string
	TextExtracted    string
	RewrittenResume  string
	UploadedByUserID *uint    `gorm:"index"`
	UploadedByUser   *User    `gorm:"foreignKey:UploadedByUserID"`
	GPTCallID        *uint    `gorm:"index"`
	GPTCall          *GPTCall `gorm:"foreignKey:GPTCallID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
