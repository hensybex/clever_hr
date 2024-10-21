// internal/mistral/mistral_init.go

package mistral

import (
	"fmt"
	"os"
	"time"

	"clever_hr_api/internal/repository"
	"github.com/go-resty/resty/v2"
)

type MistralService struct {
	apiKey       string
	url          string
	embeddingUrl string
	client       *resty.Client
	GPTCallRepo  repository.GPTCallRepository
}

func NewMistralService(gptCallRepo repository.GPTCallRepository) *MistralService {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		fmt.Println("MISTRAL_API_KEY environment variable is not set")
	}

	client := resty.New().SetTimeout(10 * time.Second)

	return &MistralService{
		apiKey:       apiKey,
		url:          "https://api.mistral.ai/v1/chat/completions",
		embeddingUrl: "https://api.mistral.ai/v1/embeddings",
		client:       client,
		GPTCallRepo:  gptCallRepo,
	}
}
