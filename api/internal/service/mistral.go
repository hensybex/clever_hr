// internal/service/mistral.go

package service

import (
	//"clever_hr/internal/dtos"
	"bufio"
	"bytes"
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	//"go_backend/internal/usecase"
	"net/http"
	"os"

	//"regexp"
	//"strings"
	"time"

	"github.com/gorilla/websocket"
	//"cuelang.org/go/cue"
	//"cuelang.org/go/cue/cuecontext"
	//"cuelang.org/go/cue/load"
)

// Define a new type for MistralModel with constants
type MistralModel string

const (
	Nemo      MistralModel = "open-mistral-nemo"
	Largest   MistralModel = "mistral-large-latest"
	Codestral MistralModel = "codestral-latest"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model          string        `json:"model"`
	ResponseFormat interface{}   `json:"response_format,omitempty"`
	Messages       []ChatMessage `json:"messages"`
	Stream         bool          `json:"stream"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type MistralService struct {
	apiKey      string
	url         string
	GPTCallRepo repository.GPTCallRepository
}

func NewMistralService(gptCallRepo repository.GPTCallRepository) *MistralService {
	apiKey := os.Getenv("MISTRAL_API_KEY")

	// Print the API key to verify it's set
	if apiKey == "" {
		fmt.Println("MISTRAL_API_KEY environment variable is not set")
	} else {
		fmt.Println("MISTRAL_API_KEY environment variable is set:", apiKey)
	}

	return &MistralService{
		apiKey:      apiKey,
		url:         "https://api.mistral.ai/v1/chat/completions",
		GPTCallRepo: gptCallRepo,
	}
}

func (ms *MistralService) CallMistral(prompt string, needJson bool, mistralModel MistralModel, entityType string, entityID uint) (string, uint, error) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := ChatRequest{
		Model:    string(mistralModel),
		Messages: messages,
	}

	if needJson {
		requestBody.ResponseFormat = map[string]string{"type": "json_object"}
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return "", 0, err
	}

	//fmt.Printf("Request body JSON: %s\n", string(jsonValue))

	var resp *http.Response
	var body []byte
	maxAttempts := 5
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		req, err := http.NewRequest("POST", ms.url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("Failed to create HTTP request: %v\n", err)
			return "", 0, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ms.apiKey)

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("HTTP request failed: %v\n", err)
			return "", 0, err
		}

		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			return "", 0, err
		}

		// Check if the response status is 429
		if resp.StatusCode == http.StatusTooManyRequests {
			//fmt.Printf("Rate limit exceeded, attempt %d of %d\n", attempts, maxAttempts)
			if attempts < maxAttempts {
				// Exponential backoff
				sleepTime := time.Duration(attempts*attempts) * time.Second
				//fmt.Printf("Sleeping for %v seconds before retrying...\n", sleepTime)
				time.Sleep(sleepTime)
				continue
			} else {
				return "", 0, fmt.Errorf("rate limit exceeded after %d attempts", maxAttempts)
			}
		}

		// Exit loop if status is not 429
		break
	}

	// Handle both JSON and non-JSON cases similarly
	responseContent := ""
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return "", 0, err
	}

	if len(chatResponse.Choices) > 0 {
		responseContent = chatResponse.Choices[0].Message.Content
	} else {
		fmt.Printf("No choices found in JSON response\n")
		return "", 0, fmt.Errorf("no choices in JSON response")
	}

	// Log the GPT call using the provided GPTCallUsecase
	gptCall := model.GPTCall{
		FinalPrompt: prompt,
		Reply:       responseContent,
		EntityType:  entityType,
		EntityID:    entityID,
	}
	gptCallId, err := ms.GPTCallRepo.CreateOne(&gptCall)
	if err != nil {
		fmt.Printf("Failed to log GPT call: %v\n", err)
		return responseContent, 0, fmt.Errorf("failed to log GPT call: %w", err)
	}

	//fmt.Printf("Successfully logged GPT call with ID: %d\n", gptCallId)

	return responseContent, gptCallId, nil
}

// Define ChatResponseChunk for streaming responses
type ChatResponseChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content,omitempty"`
		} `json:"delta"`
	} `json:"choices"`
}

func (ms *MistralService) CallMistralStream(prompt string, needJson bool, mistralModel MistralModel, entityType string, entityID uint, conn *websocket.Conn) (string, uint, error) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := ChatRequest{
		Model:    string(mistralModel),
		Messages: messages,
		Stream:   true, // Enable streaming
	}

	if needJson {
		requestBody.ResponseFormat = map[string]string{"type": "json_object"}
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return "", 0, err
	}

	var resp *http.Response
	var body []byte
	maxAttempts := 5 // Maximum number of retry attempts for rate limit
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		req, err := http.NewRequest("POST", ms.url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("Failed to create HTTP request: %v\n", err)
			return "", 0, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ms.apiKey)

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("HTTP request failed: %v\n", err)
			return "", 0, err
		}
		defer resp.Body.Close()

		// Check if the response status is 429 (rate-limited)
		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit exceeded, attempt %d of %d\n", attempts, maxAttempts)
			if attempts < maxAttempts {
				sleepTime := time.Duration(attempts*attempts) * time.Second
				fmt.Printf("Sleeping for %v seconds before retrying...\n", sleepTime)
				time.Sleep(sleepTime)
				continue
			} else {
				body, _ = io.ReadAll(resp.Body)
				fmt.Printf("Rate limit exceeded after %d attempts: %s\n", maxAttempts, string(body))
				return "", 0, fmt.Errorf("rate limit exceeded after %d attempts", maxAttempts)
			}
		}

		// Break the retry loop if the status is not 429
		if resp.StatusCode == http.StatusOK {
			break
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("Unexpected status code: %d, body: %s\n", resp.StatusCode, string(bodyBytes))
			return "", 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	// Prepare to read the streaming response
	reader := bufio.NewReader(resp.Body)
	var fullResponse string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading response: %v\n", err)
			return "", 0, err
		}

		// Handle keep-alive or empty messages
		if strings.TrimSpace(line) == "" {
			continue
		}

		// The line should be in the format: "data: {JSON}"
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimPrefix(line, "data:")
		data = strings.TrimSpace(data)

		if data == "[DONE]" {
			// Send completion status to the WebSocket client
			conn.WriteJSON(dtos.ServerMessage{Status: "Completed"})
			break
		}

		var chunk ChatResponseChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			fmt.Printf("Failed to unmarshal chunk: %v\n", err)
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content != "" {
			// Accumulate full response
			fullResponse += content

			// Send chunk to client via WebSocket
			conn.WriteJSON(dtos.ServerMessage{Result: content})
		}
	}

	// Log the GPT call
	gptCall := model.GPTCall{
		FinalPrompt: prompt,
		Reply:       fullResponse,
		EntityType:  entityType,
		EntityID:    entityID,
	}
	gptCallId, err := ms.GPTCallRepo.CreateOne(&gptCall)
	if err != nil {
		fmt.Printf("Failed to log GPT call: %v\n", err)
		return "", 0, fmt.Errorf("failed to log GPT call: %w", err)
	}

	//fmt.Printf("Successfully logged GPT call with ID: %d\n", gptCallId)

	// Close WebSocket connection after sending the completion message
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Printf("Failed to close WebSocket connection: %v\n", err)
		return "", 0, err
	}

	return fullResponse, gptCallId, nil
}
