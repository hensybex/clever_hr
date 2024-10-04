// internal/dtos/candidate_info.go

package dtos

import "clever_hr_api/internal/model"

type CandidateInfo struct {
	Candidate         model.Candidate `json:"candidate"`
	ResumePDF         string          `json:"resume_pdf"`
	ResumeID          uint            `json:"resume_id"`
	WasResumeAnalysed bool            `json:"was_resume_analysed"`
	//AnalysisResult model.ResumeAnalysisResult `json:"analysis_result"`
}
