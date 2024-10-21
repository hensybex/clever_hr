// prompts/prompts_storage/full_interview_analysis.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// FullInterviewAnalysisData contains the data for analyzing the full interview
type FullInterviewAnalysisData struct {
	InterviewMessages string // All messages from the interview
	InterviewType     string // The type of interview
}

// ToPassedData transforms FullInterviewAnalysisData into PassedData
func (f FullInterviewAnalysisData) ToPassedData() []types.PassedData {
	// Concatenate all messages

	return []types.PassedData{
		{Name: "Interview Messages", Description: "All messages from the interview", Content: f.InterviewMessages},
		{Name: "Interview Type", Description: "The type of interview", Content: f.InterviewType},
	}
}

var FullInterviewAnalysisPrompt = types.Prompt{
	BasePrompt: "You are an AI assistant analyzing a technical interview.",
	BaseTaskDesc: `**Your Task:**

1. Review the "Interview Messages" from a technical interview for the position of **"Interview Type"**.

2. Assess the candidate's performance based on their responses.

3. Identify strengths and weaknesses in the candidate's knowledge and skills.

4. Provide a recommendation on whether to proceed with the candidate.

**Respond with a JSON object containing:**

{
  "assessment": (string with detailed assessment of the candidate's performance),
  "strengths": (list of strings),
  "weaknesses": (list of strings),
  "recommendation": (string, either "proceed" or "do not proceed", with a brief explanation)
}`,
	UndefinedJSONOutputs: false,
}
