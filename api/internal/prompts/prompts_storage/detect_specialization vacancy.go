// prompts/prompts_storage/detect_specialization vacancy.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// SpecializationDetectionVacancyData contains the data for vacancy analysis and candidate extraction
type SpecializationDetectionVacancyData struct {
	VacancyDescription  string
	JobGroupType        string
	SpecializationTypes string
}

// ToPassedData transforms SpecializationDetectionVacancyData into PassedData
func (c SpecializationDetectionVacancyData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Vacancy Description", Description: "Description of a job vacancy", Content: c.VacancyDescription},
		{Name: "Job Group Type", Description: "Job Type selected by the previous assistant that fits this vacancy", Content: c.JobGroupType},
		{Name: "Specialization Types", Description: "List of possible specializations for the given Job Group Type", Content: c.SpecializationTypes},
	}
}

var SpecializationDetectionVacancyPrompt = types.Prompt{
	BasePrompt: "You're the assistant picking the right job specialization for the vacancy",
	BaseTaskDesc: `You need to analyze Vacancy Description, Job Group Type, and Specialization Types, and decide, which Specialization Type fits the given vacancy best.

Return a JSON object with the following fields:
- "Analysis": A brief analysis of the Vacancy (2-3 sentences), which reveals which Specialization Type fits this vacancy best.
- "SpecializationTypeID": Picked Job Type ID

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
