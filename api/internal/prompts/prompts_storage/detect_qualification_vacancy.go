// internal/prompts/prompts_storage/detect_qualification_vacancy.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// QualificationDetectionVacancyData contains the data for vacancy analysis and candidate extraction
type QualificationDetectionVacancyData struct {
	VacancyDescription string
	JobGroupType       string
	SpecializationType string
	QualificationTypes string
}

// ToPassedData transforms QualificationDetectionVacancyData into PassedData
func (c QualificationDetectionVacancyData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Vacancy Description", Description: "Description of a job vacancy", Content: c.VacancyDescription},
		{Name: "Job Group Type", Description: "Job Type selected by the previous assistant that fits this vacancy", Content: c.JobGroupType},
		{Name: "Specialization Type", Description: "Specialization Type selected by the previous assistant that fits this vacancy", Content: c.JobGroupType},
		{Name: "Quialification Types", Description: "List of possible quialifications", Content: c.QualificationTypes},
	}
}

var QualificationDetectionVacancyPrompt = types.Prompt{
	BasePrompt: "You're the assistant picking the right qualifications for the vacancy",
	BaseTaskDesc: `You need to analyze Vacancy Description, Job Group Type, Specialization Type and Qualification Types, and decide, which Qualification Type fits the given vacancy best.

Return a JSON object with the following fields:
- "Analysis": A brief analysis of the Vacancy (2-3 sentences), which reveals which Qualification Type fits this vacancy best.
- "QualificationTypeID": Picked Qualification Type ID

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
