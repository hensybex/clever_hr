// internal/tools/message_formatter.go

package tools

import (
	"clever_hr_api/internal/model"
	"strings"
)

func FormatMessages(messages []model.Message) string {
	var builder strings.Builder
	for _, msg := range messages {
		if msg.SentByUser {
			builder.WriteString("Human reply: ")
		} else {
			builder.WriteString("AI message: ")
		}
		builder.WriteString(msg.Message)
		builder.WriteString("\n")
	}
	return builder.String()
}
