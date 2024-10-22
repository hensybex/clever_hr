// internal/usecase/message_usecase.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
)

type MessageUsecase interface {
	CreateMessage(message *model.Message) error
	GetMessagesByTgLogin(tgLogin string) ([]model.Message, error)
}

type messageUsecase struct {
	messageRepo repository.MessageRepository
}

func NewMessageUsecase(messageRepo repository.MessageRepository) MessageUsecase {
	return &messageUsecase{messageRepo}
}

func (u *messageUsecase) CreateMessage(message *model.Message) error {
	return u.messageRepo.Create(message)
}

func (u *messageUsecase) GetMessagesByTgLogin(tgLogin string) ([]model.Message, error) {
	return u.messageRepo.GetMessagesByTgLogin(tgLogin)
}
