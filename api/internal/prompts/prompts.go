// internal/prompts/prompts.go

package prompts

import (
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/prompts/types"
)

// Prompts struct that aggregates all specific prompts
type Prompts struct {
	ResumeRewritionPrompt                types.Prompt
	FullInterviewAnalysisPrompt          types.Prompt
	CandidateExtractionPrompt            types.Prompt
	ResumeAnalysisPrompt                 types.Prompt
	GroupDetectionPrompt                 types.Prompt
	SpecializationDetectionPrompt        types.Prompt
	QualificationDetectionPrompt         types.Prompt
	GroupDetectionVacancyPrompt          types.Prompt
	SpecializationDetectionVacancyPrompt types.Prompt
	QualificationDetectionVacancyPrompt  types.Prompt
	ResumeVacancyAnalysisPrompt          types.Prompt
	ResumeStandardizationPrompt          types.Prompt
	VacancyStandardizationPrompt         types.Prompt
	InterviewMessagePrompt               types.Prompt
	FirstMessagePrompt                   types.Prompt
}

// NewPrompts initializes and returns a Prompts struct with all predefined prompts
func NewPrompts() *Prompts {
	return &Prompts{
		ResumeRewritionPrompt:                prompts_storage.ResumeRewritionPrompt,
		FullInterviewAnalysisPrompt:          prompts_storage.FullInterviewAnalysisPrompt,
		CandidateExtractionPrompt:            prompts_storage.CandidateExtractionPrompt,
		ResumeAnalysisPrompt:                 prompts_storage.ResumeAnalysisPrompt,
		GroupDetectionPrompt:                 prompts_storage.GroupDetectionPrompt,
		SpecializationDetectionPrompt:        prompts_storage.SpecializationDetectionPrompt,
		QualificationDetectionPrompt:         prompts_storage.QualificationDetectionPrompt,
		GroupDetectionVacancyPrompt:          prompts_storage.GroupDetectionVacancyPrompt,
		SpecializationDetectionVacancyPrompt: prompts_storage.SpecializationDetectionVacancyPrompt,
		QualificationDetectionVacancyPrompt:  prompts_storage.QualificationDetectionVacancyPrompt,
		ResumeVacancyAnalysisPrompt:          prompts_storage.ResumeVacancyAnalysisPrompt,
		ResumeStandardizationPrompt:          prompts_storage.ResumeStandardizationPrompt,
		VacancyStandardizationPrompt:         prompts_storage.VacancyStandardizationPrompt,
		InterviewMessagePrompt:               prompts_storage.InterviewMessagePrompt,
		FirstMessagePrompt:                   prompts_storage.FirstMessagePrompt,
	}
}
