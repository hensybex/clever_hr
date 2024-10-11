// internal/usecase/interview_result.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"errors"
)

type InterviewAnalysisResultUsecase interface {
	CreateInterviewAnalysisResult(result *model.InterviewAnalysisResult) error
	GetInterviewAnalysisResultByInterviewID(interviewID uint) (*model.InterviewAnalysisResult, error)
}

type interviewAnalysisResultUsecase struct {
	analysisRepo repository.InterviewAnalysisResultRepository
}

func NewInterviewAnalysisResultUsecase(analysisRepo repository.InterviewAnalysisResultRepository) InterviewAnalysisResultUsecase {
	return &interviewAnalysisResultUsecase{analysisRepo}
}

func (u *interviewAnalysisResultUsecase) CreateInterviewAnalysisResult(result *model.InterviewAnalysisResult) error {
	return u.analysisRepo.CreateInterviewAnalysisResult(result)
}

func (u *interviewAnalysisResultUsecase) GetInterviewAnalysisResultByInterviewID(interviewID uint) (*model.InterviewAnalysisResult, error) {
	result, err := u.analysisRepo.GetInterviewAnalysisResultByInterviewID(interviewID)
	if err != nil {
		return nil, errors.New("could not fetch interview result")
	}
	return result, nil
}
