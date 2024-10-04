// internal/usecase/resume_analysis_result_usecase.go

package usecase

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"errors"
)

type ResumeAnalysisResultUsecase interface {
	CreateResumeAnalysisResult(result *model.ResumeAnalysisResult) error
	GetResumeAnalysisResultByResumeID(resumeID uint) (*dtos.ResumeAnalysisResultDTO, error)
}

type resumeAnalysisResultUsecase struct {
	analysisRepo  repository.ResumeAnalysisResultRepository
	candidateRepo repository.CandidateRepository
	resumeRepo    repository.ResumeRepository
}

func NewResumeAnalysisResultUsecase(
	analysisRepo repository.ResumeAnalysisResultRepository,
	candidateRepo repository.CandidateRepository,
	resumeRepo repository.ResumeRepository,
) ResumeAnalysisResultUsecase {
	return &resumeAnalysisResultUsecase{analysisRepo, candidateRepo, resumeRepo}
}

func (u *resumeAnalysisResultUsecase) CreateResumeAnalysisResult(result *model.ResumeAnalysisResult) error {
	return u.analysisRepo.CreateResumeAnalysisResult(result)
}

func (u *resumeAnalysisResultUsecase) GetResumeAnalysisResultByResumeID(resumeID uint) (*dtos.ResumeAnalysisResultDTO, error) {
	result := &dtos.ResumeAnalysisResultDTO{}
	resumeAnalysisResult, err := u.analysisRepo.GetResumeAnalysisResultByResumeID(resumeID)
	if err != nil {
		return nil, errors.New("could not fetch resume analysis result")
	}
	resume, err := u.resumeRepo.GetResumeByID(resumeID)
	if err != nil {
		return nil, errors.New("could not fetch resume")
	}
	result.CandidateID = resume.CandidateID
	result.ResumeAnalysisResult = *resumeAnalysisResult
	return result, nil
}
