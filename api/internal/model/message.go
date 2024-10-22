// internal/model/message.go

package model

import (
	"time"
)

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	TgLogin    string    `json:"tg_login"`
	Message    string    `json:"message"`
	SentByUser bool      `json:"sent_by_user"`
}
