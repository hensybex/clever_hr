// internal/usecase/candidate_usecase.go

package usecase

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"fmt"
	"log"
	//"gorm.io/gorm"
)

type CandidateUsecase interface {
	CreateCandidate(candidate *model.Candidate) error
	GetCandidateByID(id uint) (*model.Candidate, error)
	GetCandidateInfo(candidateID uint) (*dtos.CandidateInfo, error)
}

type candidateUsecase struct {
	candidateRepo            repository.CandidateRepository
	resumeRepo               repository.ResumeRepository
	resumeAnalysisResultRepo repository.ResumeAnalysisResultRepository
}

func NewCandidateUsecase(
	candidateRepo repository.CandidateRepository,
	resumeRepo repository.ResumeRepository,
	resumeAnalysisResultRepo repository.ResumeAnalysisResultRepository,
) CandidateUsecase {
	return &candidateUsecase{
		candidateRepo:            candidateRepo,
		resumeRepo:               resumeRepo,
		resumeAnalysisResultRepo: resumeAnalysisResultRepo,
	}
}

func (u *candidateUsecase) CreateCandidate(candidate *model.Candidate) error {
	return u.candidateRepo.CreateCandidate(candidate)
}

func (u *candidateUsecase) GetCandidateByID(id uint) (*model.Candidate, error) {
	return u.candidateRepo.GetCandidateByID(id)
}

func (u *candidateUsecase) GetCandidateInfo(candidateID uint) (*dtos.CandidateInfo, error) {
	log.Printf("INFO: Fetching candidate details for ID: %d", candidateID)

	// Fetch candidate details
	candidate, err := u.candidateRepo.GetCandidateByID(candidateID)
	if err != nil {
		log.Printf("ERROR: Error fetching candidate with ID: %d. Error: %v", candidateID, err)
		return nil, fmt.Errorf("error fetching candidate: %v", err)
	}
	log.Printf("INFO: Successfully fetched candidate details for ID: %d", candidateID)

	// Fetch candidate's resume
	resume, err := u.resumeRepo.GetResumeByCandidateID(candidateID)
	if err != nil {
		log.Printf("ERROR: Error fetching resume for candidate ID: %d. Error: %v", candidateID, err)
		return nil, fmt.Errorf("error fetching resume: %v", err)
	}
	log.Printf("INFO: Successfully fetched resume for candidate ID: %d", candidateID)

	// Attempt to fetch analysis results for the resume, but allow a missing result
	wasResumeAnalysed, err := u.resumeAnalysisResultRepo.WasResumeAnalysed(resume.ID)
	if err != nil {
		log.Printf("ERROR: Error fetching resume analysis result check for resume ID: %d. Error: %v", resume.ID, err)
		return nil, fmt.Errorf("error fetching resume analysis result: %v", err)
	}

	// Combine the data into a single struct and return
	candidateInfo := &dtos.CandidateInfo{
		Candidate:         *candidate,
		ResumePDF:         resume.ResumePDF,
		ResumeID:          resume.ID,
		WasResumeAnalysed: wasResumeAnalysed,
	}
	log.Printf("INFO: Successfully combined candidate info, resume, and analysis result for candidate ID: %d", candidateID)

	return candidateInfo, nil
}
