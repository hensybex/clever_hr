// internal/model/interview_type.go

package model

import (
	"time"
)

type InterviewType struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
