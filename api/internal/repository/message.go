// internal/repository/message_repository.go

package repository

import (
	"clever_hr_api/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *model.Message) error
	GetMessagesByTgLogin(tgLogin string) ([]model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) Create(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetMessagesByTgLogin(tgLogin string) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("tg_login = ?", tgLogin).Order("created_at ASC").Find(&messages).Error
	return messages, err
}
