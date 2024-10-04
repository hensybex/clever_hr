// internal/prompts/prompts_storage/interview_end_check.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// InterviewEndCheckData contains the data for checking if the interview should end
type InterviewEndCheckData struct {
	FullDialogue string // The full dialogue of the interview, including the latest LLM response
}

// ToPassedData transforms InterviewEndCheckData into PassedData
func (i InterviewEndCheckData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Full Dialogue", Description: "The full dialogue of the interview, including the latest response", Content: i.FullDialogue},
	}
}

var InterviewEndCheckPrompt = types.Prompt{
	BasePrompt: "You are to determine if a technical interview should be concluded.",
	BaseTaskDesc: `**Your Task:**

1. Analyze the "Full Dialogue" of the interview.

2. Determine whether the interview should end based on the conversation.

**Instructions:**

- If all key topics have been covered, or if the candidate has significantly failed to demonstrate the required competencies, recommend ending the interview.

- If the interview should continue, respond **"false"**.

- If the interview should end, respond **"true"**.

**Respond with a single word: "true" or "false".**`,
	UndefinedJSONOutputs: false,
}
