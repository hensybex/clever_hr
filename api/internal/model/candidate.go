// internal/model/candidate.go

package model

import (
	"time"
)

type Candidate struct {
	ID               uint  `gorm:"primaryKey"`
	UploadedByUserID *uint `gorm:"index"`
	UploadedByUser   *User `gorm:"foreignKey:UploadedByUserID"`
	Name             *string
	Email            *string
	Phone            *string
	BirthDate        *time.Time
	TotalYears       *string
	PreferableJob    *string
	GPTCallID        *uint    `gorm:"index"`
	GPTCall          *GPTCall `gorm:"foreignKey:GPTCallID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
