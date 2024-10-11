package utils

import (
	"clever_hr_embeddings/internal/services"
	"context"
)

var mistralClient *mistral.MistralClient

func InitMistralClient(apiKey string) {
	mistralClient = mistral.NewMistralClient(apiKey)
}

func GenerateEmbedding(text string) ([]float64, error) {
	ctx := context.Background()
	embedding, err := mistralClient.GenerateEmbedding(ctx, text, "your-model-id")
	if err != nil {
		return nil, err
	}
	return embedding, nil
}
