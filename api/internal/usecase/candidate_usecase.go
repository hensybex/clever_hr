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
	GetCandidateInfoByTgID(tgID string) (*dtos.CandidateInfo, error)
}

type candidateUsecase struct {
	candidateRepo               repository.CandidateRepository
	resumeRepo                  repository.ResumeRepository
	userRepo                    repository.UserRepository
	resumeAnalysisResultRepo    repository.ResumeAnalysisResultRepository
	interviewRepo               repository.InterviewRepository
	interviewAnalysisResultRepo repository.InterviewAnalysisResultRepository
}

func NewCandidateUsecase(
	candidateRepo repository.CandidateRepository,
	resumeRepo repository.ResumeRepository,
	userRepo repository.UserRepository,
	resumeAnalysisResultRepo repository.ResumeAnalysisResultRepository,
	interviewRepo repository.InterviewRepository,
	interviewAnalysisResultRepo repository.InterviewAnalysisResultRepository,
) CandidateUsecase {
	return &candidateUsecase{
		candidateRepo:               candidateRepo,
		resumeRepo:                  resumeRepo,
		userRepo:                    userRepo,
		resumeAnalysisResultRepo:    resumeAnalysisResultRepo,
		interviewRepo:               interviewRepo,
		interviewAnalysisResultRepo: interviewAnalysisResultRepo,
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

func (u *candidateUsecase) GetCandidateInfoByTgID(tgID string) (*dtos.CandidateInfo, error) {
	log.Printf("INFO: Fetching candidate details for TG ID: %s", tgID)

	user, err := u.userRepo.GetUserByTgID(tgID)
	if err != nil {
		log.Printf("ERROR: Error fetching user")
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	// Fetch candidate details
	candidate, err := u.candidateRepo.GetCandidateByUploadedID(user.ID)
	if err != nil {
		log.Printf("ERROR: Error fetching candidate with TG ID: %s. Error: %v", tgID, err)
		return nil, fmt.Errorf("error fetching candidate: %v", err)
	}
	log.Printf("INFO: Successfully fetched candidate details for TG ID: %s", tgID)

	// Fetch candidate's resume
	resume, err := u.resumeRepo.GetResumeByCandidateID(candidate.ID)
	if err != nil {
		log.Printf("ERROR: Error fetching resume for candidate ID: %d. Error: %v", candidate.ID, err)
		return nil, fmt.Errorf("error fetching resume: %v", err)
	}
	log.Printf("INFO: Successfully fetched resume for candidate ID: %d", candidate.ID)

	// Attempt to fetch analysis results for the resume, but allow a missing result
	wasResumeAnalysed, err := u.resumeAnalysisResultRepo.WasResumeAnalysed(resume.ID)
	if err != nil {
		log.Printf("ERROR: Error fetching resume analysis result check for resume ID: %d. Error: %v", resume.ID, err)
		return nil, fmt.Errorf("error fetching resume analysis result: %v", err)
	}

	interview, err := u.interviewRepo.GetInterviewByResumeID(resume.ID)
	if err != nil {
		log.Printf("ERROR: Error fetching interview by resume ID: %d. Error: %v", resume.ID, err)
		return nil, fmt.Errorf("error fetching resume analysis result: %v", err)
	}

	interviewAnalysisResults, err := u.interviewAnalysisResultRepo.GetInterviewAnalysisResultByInterviewID(interview.ID)
	if err != nil {
		log.Printf("ERROR: Error fetching interview analysis results by interview ID: %d. Error: %v", interview.ID, err)
		return nil, fmt.Errorf("error fetching interview analysis result: %v", err)
	}

	// Combine the data into a single struct and return
	candidateInfo := &dtos.CandidateInfo{
		Candidate:                 *candidate,
		ResumePDF:                 resume.ResumePDF,
		ResumeID:                  resume.ID,
		WasResumeAnalysed:         wasResumeAnalysed,
		InterviewAnalysisResultID: interviewAnalysisResults.ID,
	}
	log.Printf("INFO: Successfully combined candidate info, resume, and analysis result for candidate ID: %d", candidate.ID)

	return candidateInfo, nil
}
