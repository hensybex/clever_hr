// internal/prompts/prompts_storage/detect_group_vacancy.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// GroupDetectionVacancyData contains the data for resume analysis and candidate extraction
type GroupDetectionVacancyData struct {
	VacancyDescription string
	JobTypes           string
}

// ToPassedData transforms GroupDetectionVacancyData into PassedData
func (c GroupDetectionVacancyData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Vacancy Description", Description: "Description of a job vacancy", Content: c.VacancyDescription},
		{Name: "Job Types", Description: "List of possible job types", Content: c.JobTypes},
	}
}

var GroupDetectionVacancyPrompt = types.Prompt{
	BasePrompt: "You're the assistant picking the right global job type for the vacancy, and picking appropriate title for the vacancy",
	BaseTaskDesc: `You need to analyze Vacancy Description and provided Job Types, and decide, which Job Type fits the given vacancy best. You should also pick appropriate title for the vacancy.

Return a JSON object with the following fields:
- "Analysis": A brief analysis of the Vacancy (2-3 sentences), which reveals which Job Type fits this vacancy best.
- "JobTypeID": Picked Job Type ID
- "Title": Appropriate title for the vacancy

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
