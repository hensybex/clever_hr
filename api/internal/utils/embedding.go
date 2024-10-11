// internal/utils/embedding.go

package utils

import (
	"clever_hr_api/internal/service"
	"context"
)

var mistralClient *service.MistralClient

func InitMistralClient(apiKey string) {
	mistralClient = service.NewMistralClient(apiKey)
}

func GenerateEmbedding(text string) ([]float64, error) {
	ctx := context.Background()
	embedding, err := mistralClient.GenerateEmbedding(ctx, text, "mistral-embed")
	if err != nil {
		return nil, err
	}
	return embedding, nil
}
