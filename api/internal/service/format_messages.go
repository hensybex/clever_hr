// internal/service/mistral_service.go
package service

import (
	"clever_hr_api/internal/model"
	"fmt"
)

// FormatMessages converts []model.InterviewMessage into a single string formatted as:
// User: {message[0]} LLM: {message[1]} User: {message[2]} ...
func FormatMessages(messages []model.InterviewMessage) string {
	combinedMessages := ""
	for _, message := range messages {
		if message.SentBy == model.SentByCandidate {
			combinedMessages += fmt.Sprintf("User: %s\n", message.MessageText)
		} else if message.SentBy == model.SentByLLM {
			combinedMessages += fmt.Sprintf("LLM: %s\n", message.MessageText)
		} else {
			combinedMessages += fmt.Sprintf("Employee: %s\n", message.MessageText)
		}
	}
	return combinedMessages
}
