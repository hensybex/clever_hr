package mistral

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type MistralClient struct {
	apiKey string
	client *resty.Client
}

func NewMistralClient(apiKey string) *MistralClient {
	client := resty.New()
	return &MistralClient{
		apiKey: apiKey,
		client: client,
	}
}

type EmbeddingRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"`
}

type EmbeddingResponse struct {
	ID     string         `json:"id"`
	Object string         `json:"object"`
	Model  string         `json:"model"`
	Usage  map[string]int `json:"usage"`
	Data   [][]float64    `json:"data"`
}

func (mc *MistralClient) GenerateEmbedding(ctx context.Context, text string, model string) ([]float64, error) {
	reqBody := EmbeddingRequest{
		Input:          text,
		Model:          model,
		EncodingFormat: "float",
	}

	var respBody EmbeddingResponse

	resp, err := mc.client.R().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", mc.apiKey)).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&respBody).
		Post("https://api.mistral.ai/v1/embeddings")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("failed to generate embedding")
	}

	if len(respBody.Data) == 0 || len(respBody.Data[0]) == 0 {
		return nil, errors.New("empty embedding data")
	}

	// Assuming we're interested in the first embedding
	embedding := respBody.Data[0]
	return embedding, nil
}
