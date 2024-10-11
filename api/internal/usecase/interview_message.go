// internal/usecase/interview_message.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"errors"
)

type InterviewMessageUsecase interface {
	CreateMessage(message *model.InterviewMessage) error
	GetMessagesByInterviewID(interviewID uint) ([]model.InterviewMessage, error)
}

type interviewMessageUsecase struct {
	messageRepo repository.InterviewMessageRepository
}

func NewInterviewMessageUsecase(messageRepo repository.InterviewMessageRepository) InterviewMessageUsecase {
	return &interviewMessageUsecase{messageRepo}
}

func (u *interviewMessageUsecase) CreateMessage(message *model.InterviewMessage) error {
	return u.messageRepo.CreateMessage(message)
}

func (u *interviewMessageUsecase) GetMessagesByInterviewID(interviewID uint) ([]model.InterviewMessage, error) {
	messages, err := u.messageRepo.GetMessagesByInterviewID(interviewID)
	if err != nil {
		return nil, errors.New("could not fetch interview messages")
	}
	return messages, nil
}
