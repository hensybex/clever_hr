// internal/model/interview_message.go

package model

import (
	"time"
)

type SentBy string

const (
	SentByCandidate SentBy = "candidate"
	SentByLLM       SentBy = "llm"
)

type InterviewMessage struct {
	ID          uint      `gorm:"primaryKey"`
	InterviewID uint      `gorm:"index"`
	Interview   Interview `gorm:"foreignKey:InterviewID"`
	MessageText string    `gorm:"type:text"`
	SentBy      SentBy    `gorm:"type:varchar(20)"`
	GPTCallID   *uint     `gorm:"index"`
	GPTCall     *GPTCall  `gorm:"foreignKey:GPTCallID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
