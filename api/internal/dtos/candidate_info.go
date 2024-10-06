// internal/dtos/candidate_info.go

package dtos

import "clever_hr_api/internal/model"

type CandidateInfo struct {
	Candidate                 model.Candidate `json:"candidate"`
	ResumePDF                 string          `json:"resume_pdf"`
	ResumeID                  uint            `json:"resume_id"`
	WasResumeAnalysed         bool            `json:"was_resume_analysed"`
	InterviewAnalysisResultID uint            `json:"interview_analysis_result_id"`
	//AnalysisResult model.ResumeAnalysisResult `json:"analysis_result"`
}
