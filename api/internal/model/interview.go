// internal/model/interview.go

package model

import (
	"time"
)

type InterviewStatus string

const (
	Ongoing            InterviewStatus = "ongoing"
	InterviewFinished  InterviewStatus = "finished"
	WaitingForAnalysis InterviewStatus = "waiting_for_analysis"
)

type Interview struct {
	ID              uint            `gorm:"primaryKey"`
	ResumeID        uint            `gorm:"index"`
	Resume          Resume          `gorm:"foreignKey:ResumeID"`
	InterviewTypeID uint            `gorm:"index"`
	InterviewType   InterviewType   `gorm:"foreignKey:InterviewTypeID"`
	Status          InterviewStatus `gorm:"type:varchar(20)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
