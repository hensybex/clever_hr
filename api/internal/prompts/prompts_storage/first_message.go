// internal/prompts/prompts_storage/detect_group_vacancy.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// FirstMessageData contains the data for resume analysis and candidate extraction
type FirstMessageData struct {
	AnalysisData string
}

// ToPassedData transforms FirstMessageData into PassedData
func (c FirstMessageData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Analysis Data", Description: "Results of analysis of potential job candidate", Content: c.AnalysisData},
	}
}

var FirstMessagePrompt = types.Prompt{
	BasePrompt: "You're the assistant that writes a starting message for the conversation with a potential job candidate, based on Analysis Data, aiming to send candidate conversation starting message, and outline a list of questions, which might unveal some revealed characteristics of candidate based on the analysis.",
	BaseTaskDesc: `You need to analyze Analysis Data and a proper conversation starting message. 

Return a JSON object with the following field:
- "StartMessage": The most suitable starting message aiming to extend the initial analysis.

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
