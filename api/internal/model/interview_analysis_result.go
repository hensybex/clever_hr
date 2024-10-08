// internal/model/interview_analysis_result.go

package model

import (
	"time"
)

type ResultStatus string

const (
	BeingAnalysed  ResultStatus = "being_analysed"
	ResultFinished ResultStatus = "finished"
)

type InterviewAnalysisResult struct {
	ID             uint         `gorm:"primaryKey"`
	InterviewID    uint         `gorm:"index"`
	Interview      Interview    `gorm:"foreignKey:InterviewID"`
	GPTCallID      *uint        `gorm:"index"`
	GPTCall        *GPTCall     `gorm:"foreignKey:GPTCallID"`
	ResultStatus   ResultStatus `gorm:"type:varchar(20)"`
	Assessment     string       `gorm:"type:text"`
	Strengths      string       `gorm:"type:text"`
	Weaknesses     string       `gorm:"type:text"`
	Recommendation string       `gorm:"type:text"`
	Reason         string       `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
