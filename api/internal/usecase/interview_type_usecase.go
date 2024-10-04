// internal/usecase/interview_type_usecase.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
)

type InterviewTypeUsecase interface {
	ListInterviewTypes() ([]model.InterviewType, error)
}

type interviewTypeUsecase struct {
	repo repository.InterviewTypeRepository
}

func NewInterviewTypeUsecase(repo repository.InterviewTypeRepository) InterviewTypeUsecase {
	return &interviewTypeUsecase{repo}
}

// ListInterviewTypes returns all interview types
func (u *interviewTypeUsecase) ListInterviewTypes() ([]model.InterviewType, error) {
	return u.repo.ListInterviewTypes()
}
