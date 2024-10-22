// internal/handlers/message_handler.go

package handlers

import (
	"clever_hr_api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageUsecase usecase.MessageUsecase
}

func NewMessageHandler(messageUsecase usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
	}
}

// GetMessagesByTgLogin handles GET /api/messages/get_by_tg_login/:tg_login
func (h *MessageHandler) GetMessagesByTgLogin(c *gin.Context) {
	tgLogin := c.Param("tg_login")
	messages, err := h.messageUsecase.GetMessagesByTgLogin(tgLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}
