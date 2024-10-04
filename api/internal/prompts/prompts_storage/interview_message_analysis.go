// internal/prompts/prompts_storage/interview_message_analysis.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// InterviewMessageAnalysisData contains the data for analyzing an interview message
type InterviewMessageAnalysisData struct {
	PreviousMessages string
	LatestMessage    string // The latest message from the candidate
	InterviewType    string // The type of interview (e.g., "Python Junior", "Flutter Middle")
}

// ToPassedData transforms InterviewMessageAnalysisData into PassedData
func (i InterviewMessageAnalysisData) ToPassedData() []types.PassedData {
	// Concatenate all previous messages
	return []types.PassedData{
		{Name: "Previous Messages", Description: "All previous messages in the interview", Content: i.PreviousMessages},
		{Name: "Latest Message", Description: "The latest message from the candidate", Content: i.LatestMessage},
		{Name: "Interview Type", Description: "The type of interview", Content: i.InterviewType},
	}
}

var InterviewMessageAnalysisPrompt = types.Prompt{
	BasePrompt: "You are acting as an interviewer conducting a technical interview.",
	BaseTaskDesc: `**Your Task:**

1. Review the "Previous Messages" and the "Latest Message" from the candidate.

2. As an interviewer for the position of **"Interview Type"**, generate an appropriate response or question to continue the interview, focusing on assessing the candidate's skills relevant to the position.

3. Ensure your response is professional and maintains the flow of the interview.

**Instructions:**

- Keep the interview focused on the technical skills required for the position.

- If the candidate's response indicates a misunderstanding, politely clarify or rephrase your question.

- Do not reveal the correct answers; instead, guide the candidate to think critically.

**Respond with the next message in the interview.**`,
	UndefinedJSONOutputs: false,
}
