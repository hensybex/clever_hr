// mistral/mistral_embeddings.go

package mistral

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

// Method for generating embeddings
func (ms *MistralService) GenerateEmbedding(text string) ([]float64, error) {
	reqBody := map[string]interface{}{
		"input":           text,
		"model":           "mistral-embed",
		"encoding_format": "float",
	}

	jsonValue, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Set up retry mechanism with exponential backoff
	maxAttempts := 5
	var resp *resty.Response

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		// Make the HTTP request
		resp, err = ms.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+ms.apiKey).
			SetBody(jsonValue).
			Post(ms.embeddingUrl)

		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}

		// If we get a 429 (rate limit), handle retries with exponential backoff
		if resp.StatusCode() == http.StatusTooManyRequests {
			if attempts < maxAttempts {
				// Exponential backoff
				sleepTime := time.Duration(attempts*attempts) * time.Second
				log.Printf("Rate limit exceeded, attempt %d of %d. Retrying in %v seconds...", attempts, maxAttempts, sleepTime)
				time.Sleep(sleepTime)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d attempts", maxAttempts)
			}
		}

		// Break out of loop if status is not 429
		break
	}

	// Check if the response status code is not 200
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Parse the response
	var respBody EmbeddingResponse
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	// Check if the data contains at least one embedding
	if len(respBody.Data) == 0 || len(respBody.Data[0].Embedding) == 0 {
		return nil, errors.New("empty embedding data")
	}

	// Return the first embedding in the data array
	return respBody.Data[0].Embedding, nil
}
