// internal/prompts/prompts.go

package prompts

import (
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/prompts/types"
)

// Prompts struct that aggregates all specific prompts
type Prompts struct {
	ResumeRewritionPrompt          types.Prompt
	InterviewMessageAnalysisPrompt types.Prompt
	InterviewEndCheckPrompt        types.Prompt
	FullInterviewAnalysisPrompt    types.Prompt
	CandidateExtractionPrompt      types.Prompt
	ResumeAnalysisPrompt           types.Prompt
}

// NewPrompts initializes and returns a Prompts struct with all predefined prompts
func NewPrompts() *Prompts {
	return &Prompts{
		ResumeRewritionPrompt:          prompts_storage.ResumeRewritionPrompt,
		InterviewMessageAnalysisPrompt: prompts_storage.InterviewMessageAnalysisPrompt,
		InterviewEndCheckPrompt:        prompts_storage.InterviewEndCheckPrompt,
		FullInterviewAnalysisPrompt:    prompts_storage.FullInterviewAnalysisPrompt,
		CandidateExtractionPrompt:      prompts_storage.CandidateExtractionPrompt,
		ResumeAnalysisPrompt:           prompts_storage.ResumeAnalysisPrompt,
	}
}
