// internal/prompts/prompts_storage/detect_qualification.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// QualificationDetectionData contains the data for resume analysis and candidate extraction
type QualificationDetectionData struct {
	ResumeText         string
	JobGroupType       string
	SpecializationType string
	QualificationTypes string
}

// ToPassedData transforms QualificationDetectionData into PassedData
func (c QualificationDetectionData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "The full text extracted from the candidate's resume", Content: c.ResumeText},
		{Name: "Job Group Type", Description: "Job Type selected by the previous assistant that fits this resume", Content: c.JobGroupType},
		{Name: "Specialization Type", Description: "Specialization Type selected by the previous assistant that fits this resume", Content: c.JobGroupType},
		{Name: "Quialification Types", Description: "List of possible quialifications", Content: c.QualificationTypes},
	}
}

var QualificationDetectionPrompt = types.Prompt{
	BasePrompt: "You're the assistant picking the right qualifications for the resume",
	BaseTaskDesc: `You need to analyze Resume Text, Job Group Type, Specialization Type and Qualification Types, and decide, which Qualification Type fits the given resume best.

Return a JSON object with the following fields:
- "Analysis": A brief analysis of the resume (2-3 sentences), which reveals which Qualification Type fits this resume best.
- "QualificationTypeID": Picked Qualification Type ID

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
