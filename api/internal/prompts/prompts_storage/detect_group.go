// internal/prompts/prompts_storage/detect_group.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// GroupDetectionData contains the data for resume analysis and candidate extraction
type GroupDetectionData struct {
	ResumeText string
	JobTypes   string
}

// ToPassedData transforms GroupDetectionData into PassedData
func (c GroupDetectionData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "The full text extracted from the candidate's resume", Content: c.ResumeText},
		{Name: "Job Types", Description: "List of possible jobs", Content: c.JobTypes},
	}
}

var GroupDetectionPrompt = types.Prompt{
	BasePrompt: "You're the assistant picking the right global job type for the resume, and retrieving full name from the resume",
	BaseTaskDesc: `You need to analyze Resume Text and provided Job Types, and decide, which Job Type fits the given resume best. You should also find full name from the resume.

Return a JSON object with the following fields:
- "Analysis": A brief analysis of the Resume (2-3 sentences), which reveals which Job Type fits this resume text best.
- "JobTypeID": Picked Job Type ID
- "FullName": Full name of the person found in the resume

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
