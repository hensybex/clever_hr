// internal/prompts/prompts_storage/detect_group_vacancy.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// InterviewMessageData contains the data for resume analysis and candidate extraction
type InterviewMessageData struct {
	PreviousDialog string
}

// ToPassedData transforms InterviewMessageData into PassedData
func (c InterviewMessageData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Previous dialog", Description: "Dialog between human and LLM", Content: c.PreviousDialog},
	}
}

var InterviewMessagePrompt = types.Prompt{
	BasePrompt: "You're the assistant understanding the previous dialog and generating a suitable LLM reply in Russian for the last human's message",
	BaseTaskDesc: `You need to analyze Previous Dialog and decide, which reply to provide to user. 

Return a JSON object with the following field:
- "Response": The most suitable response for the last human's message based on the dialogue provided in Russian language.

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
