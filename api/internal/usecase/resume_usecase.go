// internal/usecase/resume_usecase.go

package usecase

import (
	"bytes"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/prompts"
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type ResumeUsecase interface {
	UploadResume(resume *model.Resume, filePath string) error
	RunResumeAnalysis(resumeID uint) (*model.ResumeAnalysisResult, error)
	GetResumeAnalysisResult(resumeID uint) (*model.ResumeAnalysisResult, error)
}

type resumeUsecase struct {
	resumeRepo     repository.ResumeRepository
	analysisRepo   repository.ResumeAnalysisResultRepository
	candidateRepo  repository.CandidateRepository
	userRepo       repository.UserRepository
	mistralService service.MistralService
}

func NewResumeUsecase(
	resumeRepo repository.ResumeRepository,
	analysisRepo repository.ResumeAnalysisResultRepository,
	candidateRepo repository.CandidateRepository,
	userRepo repository.UserRepository,
	mistralService service.MistralService,
) ResumeUsecase {
	return &resumeUsecase{resumeRepo, analysisRepo, candidateRepo, userRepo, mistralService}
}

type CandidateData struct {
	FullName        *string `json:"FullName"`
	Email           *string `json:"Email"`
	Phone           *string `json:"Phone"`
	BirthDate       *string `json:"BirthDate"`
	TotalYears      *string `json:"TotalYears"`
	PreferableJob   *string `json:"PreferableJob"`
	RewrittenResume string  `json:"RewrittenResume"`
}

func (u *resumeUsecase) UploadResume(resume *model.Resume, filePath string) error {
	// Log the start of the function and input parameters
	log.Printf("Starting EmployeeUploadResume. filePath: %s, resume ID: %v, uploaded by user ID: %v", filePath, resume.ID, resume.UploadedByUserID)

	// Extract text from PDF
	extractedText, err := u.extractTextFromPDF(filePath)
	if err != nil {
		log.Printf("Error extracting text from PDF. filePath: %s, error: %v", filePath, err)
		return err
	}
	log.Printf("Text successfully extracted from PDF. filePath: %s, extracted text length: %d", filePath, len(extractedText))

	// Construct LLM prompt
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	dialogData := prompts_storage.CandidateExtractionData{ResumeText: extractedText}
	prompt, err := pc.GetPrompt(allPrompts.CandidateExtractionPrompt, dialogData, "", true)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("LLM prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the resume
	llmResponse, gptCallID, err := u.mistralService.CallMistral(prompt, true, service.Nemo, "ResumeAnalysis", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var candidateData CandidateData
	if err := json.Unmarshal([]byte(llmResponse), &candidateData); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Parse birth date if provided
	var birthDate *time.Time
	if candidateData.BirthDate != nil && *candidateData.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *candidateData.BirthDate)
		if err == nil {
			birthDate = &parsedDate
			log.Printf("Birthdate parsed successfully. Birthdate: %s", parsedDate.Format("2006-01-02"))
		} else {
			log.Printf("Error parsing birthdate. raw birthdate: %s, error: %v", *candidateData.BirthDate, err)
		}
	}

	// Parse total years of experience
	totalYearsStr := candidateData.TotalYears

	// Get user by Telegram ID
	var telegramID string
	var user *model.User
	if resume.UploadedByUserID != nil {
		telegramID = strconv.FormatUint(uint64(*resume.UploadedByUserID), 10)
		user, err = u.userRepo.GetUserByID(*resume.UploadedByUserID)
		if err != nil {
			log.Printf("Error getting user by Telegram ID. telegramID: %s, error: %v", telegramID, err)
			return fmt.Errorf("error getting user: %v", err)
		}
	} else {
		telegramID = strconv.FormatUint(uint64(resume.CandidateID), 10)
		user, err = u.userRepo.GetUserByID(resume.CandidateID)
		if err != nil {
			log.Printf("Error getting user by Telegram ID. telegramID: %s, error: %v", telegramID, err)
			return fmt.Errorf("error getting user: %v", err)
		}
	}
	log.Printf("Fetching user by Telegram ID: %s", telegramID)

	log.Printf("User fetched successfully. User ID: %v", user.ID)

	// Create the candidate in the database
	candidate := &model.Candidate{
		UploadedByUserID: &user.ID,
		Name:             candidateData.FullName,
		Email:            candidateData.Email,
		Phone:            candidateData.Phone,
		BirthDate:        birthDate,
		TotalYears:       totalYearsStr,
		PreferableJob:    candidateData.PreferableJob,
		GPTCallID:        &gptCallID,
	}
	if err := u.candidateRepo.CreateCandidate(candidate); err != nil {
		log.Printf("Error creating candidate. error: %v", err)
		return fmt.Errorf("error creating candidate: %v", err)
	}
	log.Printf("Candidate created successfully. Candidate ID: %v", candidate.ID)

	// Save the extracted text and rewritten resume in the database
	resume.CandidateID = candidate.ID
	resume.TextExtracted = extractedText
	resume.RewrittenResume = candidateData.RewrittenResume
	resume.UploadedByUserID = &user.ID
	log.Printf("Saving resume to the database. Resume ID: %v", resume.ID)
	if err := u.resumeRepo.CreateResume(resume); err != nil {
		log.Printf("Error saving resume. error: %v", err)
		return fmt.Errorf("error saving resume: %v", err)
	}
	log.Printf("Resume saved successfully. Resume ID: %v", resume.ID)

	// Log the successful completion of the function
	log.Printf("EmployeeUploadResume completed successfully for Resume ID: %v", resume.ID)

	return nil
}

// Extract text from the PDF file using pdftotext
func (u *resumeUsecase) extractTextFromPDF(filePath string) (string, error) {
	// Command and arguments
	cmd := exec.CommandContext(context.Background(), "pdftotext", "-layout", "-nopgbrk", filePath, "-")

	// Capture the command's output and error
	var buf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &stderrBuf

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		// Log the error and any output from stderr
		return "", fmt.Errorf("error running pdftotext: %v, stderr: %s", err, stderrBuf.String())
	}

	// Return the extracted text
	return buf.String(), nil
}

var resumeAnalysisResponseStruct struct {
	ProfessionalSummaryAndCareerNarrative          model.AnalysisField `json:"ProfessionalSummaryAndCareerNarrative"`
	WorkExperienceAndImpact                        model.AnalysisField `json:"WorkExperienceAndImpact"`
	EducationAndContinuousLearning                 model.AnalysisField `json:"EducationAndContinuousLearning"`
	SkillsAndTechnologicalProficiency              model.AnalysisField `json:"SkillsAndTechnologicalProficiency"`
	SoftSkillsAndEmotionalIntelligence             model.AnalysisField `json:"SoftSkillsAndEmotionalIntelligence"`
	LeadershipInnovationAndProblemSolvingPotential model.AnalysisField `json:"LeadershipInnovationAndProblemSolvingPotential"`
	CulturalFitAndValueAlignment                   model.AnalysisField `json:"CulturalFitAndValueAlignment"`
	AdaptabilityResilienceAndWorkEthic             model.AnalysisField `json:"AdaptabilityResilienceAndWorkEthic"`
	LanguageProficiencyAndCommunication            model.AnalysisField `json:"LanguageProficiencyAndCommunicationSkills"`
	ProfessionalAffiliationsAndCommunity           model.AnalysisField `json:"ProfessionalAffiliationsAndCommunityEngagement"`
}

func (u *resumeUsecase) RunResumeAnalysis(resumeID uint) (*model.ResumeAnalysisResult, error) {
	// Fetch the resume data
	resume, err := u.resumeRepo.GetResumeByID(resumeID)
	if err != nil {
		return nil, errors.New("resume not found")
	}

	// Construct LLM prompt
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	resumeAnalysisData := prompts_storage.ResumeAnalysisData{ResumeText: resume.RewrittenResume}
	prompt, err := pc.GetPrompt(allPrompts.ResumeAnalysisPrompt, resumeAnalysisData, "Russian", true)
	if err != nil {
		return nil, errors.New("error generating LLM prompt for resume analysis")
	}

	// Call LLM to rewrite the resume in a structured format
	resumeAnalysis, gptCallID, err := u.mistralService.CallMistral(prompt, true, service.Nemo, "ResumeAnalysis", 0)
	if err != nil {
		return nil, errors.New("error calling LLM for resume analysis")
	}

	// Parse LLM response into the structured data
	if err := json.Unmarshal([]byte(resumeAnalysis), &resumeAnalysisResponseStruct); err != nil {
		return nil, errors.New("error parsing LLM response")
	}

	// Create the analysis result struct
	analysisResult := &model.ResumeAnalysisResult{
		ResumeID:       resumeID,
		AnalysisStatus: model.AnalysisFinished,

		// Store the parsed analysis fields
		ProfessionalSummaryAndCareerNarrative:          resumeAnalysisResponseStruct.ProfessionalSummaryAndCareerNarrative,
		WorkExperienceAndImpact:                        resumeAnalysisResponseStruct.WorkExperienceAndImpact,
		EducationAndContinuousLearning:                 resumeAnalysisResponseStruct.EducationAndContinuousLearning,
		SkillsAndTechnologicalProficiency:              resumeAnalysisResponseStruct.SkillsAndTechnologicalProficiency,
		SoftSkillsAndEmotionalIntelligence:             resumeAnalysisResponseStruct.SoftSkillsAndEmotionalIntelligence,
		LeadershipInnovationAndProblemSolvingPotential: resumeAnalysisResponseStruct.LeadershipInnovationAndProblemSolvingPotential,
		CulturalFitAndValueAlignment:                   resumeAnalysisResponseStruct.CulturalFitAndValueAlignment,
		AdaptabilityResilienceAndWorkEthic:             resumeAnalysisResponseStruct.AdaptabilityResilienceAndWorkEthic,
		LanguageProficiencyAndCommunication:            resumeAnalysisResponseStruct.LanguageProficiencyAndCommunication,
		ProfessionalAffiliationsAndCommunity:           resumeAnalysisResponseStruct.ProfessionalAffiliationsAndCommunity,

		// Link GPT call ID for tracking
		GPTCallID: &gptCallID,
	}

	// Save the analysis result in the database
	if err := u.analysisRepo.CreateResumeAnalysisResult(analysisResult); err != nil {
		return nil, err
	}

	return analysisResult, nil
}

func (u *resumeUsecase) GetResumeAnalysisResult(resumeID uint) (*model.ResumeAnalysisResult, error) {
	return u.analysisRepo.GetResumeAnalysisResultByResumeID(resumeID)
}
