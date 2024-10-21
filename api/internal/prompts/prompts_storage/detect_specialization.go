// internal/prompts/prompts_storage/detect_specialization.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// SpecializationDetectionData contains the data for resume analysis and candidate extraction
type SpecializationDetectionData struct {
	ResumeText          string
	JobGroupType        string
	SpecializationTypes string
}

// ToPassedData transforms SpecializationDetectionData into PassedData
func (c SpecializationDetectionData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "The full text extracted from the candidate's resume", Content: c.ResumeText},
		{Name: "Job Group Type", Description: "Job Type selected by the previous assistant that fits this resume", Content: c.JobGroupType},
		{Name: "Specialization Types", Description: "List of possible specializations for the given Job Group Type", Content: c.SpecializationTypes},
	}
}

var SpecializationDetectionPrompt = types.Prompt{
	BasePrompt: "You're the assistant picking the right job specialization for the resume",
	BaseTaskDesc: `You need to analyze Resume Text, Job Group Type, and Specialization Types, and decide, which Specialization Type fits the given resume best.

Return a JSON object with the following fields:
- "Analysis": A brief analysis of the Resume (2-3 sentences), which reveals which Specialization Type fits this resume best.
- "SpecializationTypeID": Picked Job Type ID

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
