// dtos/resume_analysis_result_dto.go

package dtos

import "clever_hr_api/internal/model"

type ResumeAnalysisResultDTO struct {
	ResumeAnalysisResult model.ResumeAnalysisResult `json:"resume_analysis_result"`
	CandidateID          uint                       `json:"candidate_id"`
}
