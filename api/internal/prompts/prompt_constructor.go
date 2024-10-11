// internal/prompts/prompt_constructor.go

package prompts

import (
	"clever_hr_api/internal/prompts/types"
	"fmt"
)

// PromptConstructor is used to construct prompts with specific content
type PromptConstructor struct{}

// NewPromptConstructor initializes and returns a PromptConstructor
func NewPromptConstructor() *PromptConstructor {
	return &PromptConstructor{}
}

// GetPrompt constructs a prompt based on the provided prompt template and data
func (pc *PromptConstructor) GetPrompt(prompt types.Prompt, data types.PromptData, language string, repeatLanguage bool) (string, error) {
	prompt.Language = language
	passedData := data.ToPassedData()

	// Construct the list of passed data
	passedDataList := "You will receive:\n"
	passedDataContentStr := ""
	for _, data := range passedData {
		passedDataList += fmt.Sprintf("%s - %s\n", data.Name, data.Description)
		passedDataContentStr += fmt.Sprintf("%s - %s\n\n", data.Name, data.Content)
	}

	// Construct the language instruction if applicable
	languageInstruction := ""
	if prompt.Language != "" {
		languageInstruction = fmt.Sprintf("Fully respond in %s.\n", prompt.Language)
	}

	// Construct the JSON instruction if applicable
	jsonInstruction := ""
	if len(prompt.JSONStruct) > 0 {
		jsonInstruction = "Your response should be a structured JSON with the following keys:\n"
		for _, js := range prompt.JSONStruct {
			jsonInstruction += fmt.Sprintf("%s: %s\n", js.Key, js.Description)
		}
	}

	if prompt.UndefinedJSONOutputs {
		jsonInstruction += "(You should decide the number of output probabilities yourself)\n"
	}

	// Construct the final prompt with optional language instruction
	finalPrompt := ""
	if languageInstruction != "" {
		finalPrompt = languageInstruction // Always include at the start
	}

	finalPrompt += prompt.BasePrompt + "\n\n" + passedDataList + prompt.BaseTaskDesc + "\n\n" + passedDataContentStr

	if languageInstruction != "" && repeatLanguage {
		finalPrompt += languageInstruction // Add again at the end if repeatLanguage is true
	}

	// Append JSON structure instruction if available
	finalPrompt += jsonInstruction

	// Append full examples of output JSONs if available
	if len(prompt.JSONExamples) > 0 {
		finalPrompt += "\n\nFull examples of output JSONs:\n"
		for _, example := range prompt.JSONExamples {
			finalPrompt += fmt.Sprintf("Request: %s\nResponse:\n%s\n\n", example.Request, example.Response)
		}
	}

	return finalPrompt, nil
}
