// internal/repository/interview_message_repository.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type InterviewMessageRepository interface {
	CreateMessage(message *model.InterviewMessage) error
	GetMessagesByInterviewID(interviewID uint) ([]model.InterviewMessage, error)
}

type interviewMessageRepository struct {
	db *gorm.DB
}

func NewInterviewMessageRepository(db *gorm.DB) InterviewMessageRepository {
	return &interviewMessageRepository{db}
}

func (r *interviewMessageRepository) CreateMessage(message *model.InterviewMessage) error {
	return r.db.Create(message).Error
}

func (r *interviewMessageRepository) GetMessagesByInterviewID(interviewID uint) ([]model.InterviewMessage, error) {
	var messages []model.InterviewMessage
	err := r.db.Where("interview_id = ?", interviewID).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
