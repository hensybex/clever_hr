// internal/usecase/resume.go

package usecase

import (
	"bytes"
	"clever_hr_api/internal/mistral"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/prompts"
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type ResumeUsecase interface {
	GetResumeByID(resumeID uint) (*model.Resume, error)
	UploadResume(resume *model.Resume, filePath string) error
	RunResumeAnalysis(resumeID uint) (*model.ResumeAnalysisResult, error)
	GetResumeAnalysisResult(resumeID uint) (*model.ResumeAnalysisResult, error)
	UpdateAllResumeEmbeddings() error
}

type resumeUsecase struct {
	resumeRepo         repository.ResumeRepository
	embeddingRepo      repository.EmbeddingRepository
	analysisRepo       repository.ResumeAnalysisResultRepository
	userRepo           repository.UserRepository
	jobGroupRepo       repository.JobGroupRepository
	specializationRepo repository.SpecializationRepository
	qualificationRepo  repository.QualificationRepository
	mistralService     mistral.MistralService
}

func NewResumeUsecase(
	resumeRepo repository.ResumeRepository,
	embeddingRepo repository.EmbeddingRepository,
	analysisRepo repository.ResumeAnalysisResultRepository,
	userRepo repository.UserRepository,
	jobGroupRepo repository.JobGroupRepository,
	specializationRepo repository.SpecializationRepository,
	qualificationRepo repository.QualificationRepository,
	mistralService mistral.MistralService,
) ResumeUsecase {
	return &resumeUsecase{
		resumeRepo,
		embeddingRepo,
		analysisRepo,
		userRepo,
		jobGroupRepo,
		specializationRepo,
		qualificationRepo,
		mistralService,
	}
}

func (u *resumeUsecase) GetResumeByID(resumeID uint) (*model.Resume, error) {
	return u.resumeRepo.GetByID(resumeID)
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

type GroupDetectionResponse struct {
	Analysis  string `json:"Analysis"`
	JobTypeID uint   `json:"JobTypeID"`
	FullName  string `json:"FullName"`
}

type SpecializationDetectionResponse struct {
	Analysis             string `json:"Analysis"`
	SpecializationTypeID uint   `json:"SpecializationTypeID"`
}

type QualificationDetectionResponse struct {
	Analysis            string `json:"Analysis"`
	QualificationTypeID uint   `json:"QualificationTypeID"`
}

func cleanResumeText(input string) string {
	// Step 1: Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")

	// Step 2: Remove spaces around punctuation
	cleaned = strings.ReplaceAll(cleaned, " ,", ",")
	cleaned = strings.ReplaceAll(cleaned, " .", ".")
	cleaned = strings.ReplaceAll(cleaned, " :", ":")
	cleaned = strings.ReplaceAll(cleaned, " ;", ";")

	// Step 3: Trim leading and trailing spaces
	cleaned = strings.TrimSpace(cleaned)

	// Step 4: (Optional) You can add more formatting logic here as needed

	return cleaned
}

func (u *resumeUsecase) UploadResume(resume *model.Resume, filePath string) error {
	// Log the start of the function and input parameters
	log.Printf("Starting UploadResume. filePath: %s", filePath)

	// Extract text from PDF
	extractedText, err := u.extractTextFromPDF(filePath)
	if err != nil {
		log.Printf("Error extracting text from PDF. filePath: %s, error: %v", filePath, err)
		return err
	}
	log.Printf("Text successfully extracted from PDF. filePath: %s, extracted text length: %d", filePath, len(extractedText))

	// Clean and format the resume text
	formattedText := cleanResumeText(extractedText)

	// Get Groups Data
	jobGroups, err := u.jobGroupRepo.GetAll()
	if err != nil {
		log.Printf("Error getting all job groups. error: %v", err)
		return err
	}
	var jobGroupsString string

	for _, group := range jobGroups {
		jobGroupsString += fmt.Sprintf("Group name: %s, id: %d; ", group.TitleEn, group.ID)
	}

	jobGroupsString = strings.TrimSuffix(jobGroupsString, "; ")

	// Step 0: Construct LLM prompt for Resume Standarization
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	standarizationData := prompts_storage.ResumeStandardizationData{ResumeText: formattedText}
	prompt, err := pc.GetPrompt(allPrompts.ResumeStandardizationPrompt, standarizationData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("Resume standarization prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the resume
	llmResponse, _, err := u.mistralService.CallMistral(prompt, false, mistral.Largest, "ResumeStandarization", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}
	log.Printf("Resume standarization performed successfully successfully.")

	standarizedResumeText := llmResponse

	// Step 1: Construct LLM prompt for Group Detection
	groupDetectionData := prompts_storage.GroupDetectionData{ResumeText: standarizedResumeText, JobTypes: jobGroupsString}
	prompt, err = pc.GetPrompt(allPrompts.GroupDetectionPrompt, groupDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("GroupDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the resume
	llmResponse, _, err = u.mistralService.CallMistral(prompt, true, mistral.Nemo, "GroupDetection", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var groupDetectionResponse GroupDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &groupDetectionResponse); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Step 2: Construct LLM prompt for Specialization Detection
	// Obtain Job Group
	jobGroup, err := u.jobGroupRepo.GetByID(int(groupDetectionResponse.JobTypeID))
	if err != nil {
		log.Printf("Error getting job group by ID. error: %v", err)
		return err
	}

	// Obtain correct specialization
	specializations, err := u.specializationRepo.GetAllByGroupID(int(groupDetectionResponse.JobTypeID))
	if err != nil {
		log.Printf("Error getting all job groups. error: %v", err)
		return err
	}
	var specializationsString string

	for _, specialization := range specializations {
		specializationsString += fmt.Sprintf("Specialization name: %s, id: %d; ", specialization.TitleEn, specialization.ID)
	}

	specializationsString = strings.TrimSuffix(specializationsString, "; ")

	specializationDetectionData := prompts_storage.SpecializationDetectionData{ResumeText: standarizedResumeText, JobGroupType: jobGroup.TitleEn, SpecializationTypes: specializationsString}
	specializationPrompt, err := pc.GetPrompt(allPrompts.SpecializationDetectionPrompt, specializationDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("SpecializationDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the resume
	llmResponse, _, err = u.mistralService.CallMistral(specializationPrompt, true, mistral.Nemo, "SpecializationDetection", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var specializationResponse SpecializationDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &specializationResponse); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Step 3: Construct LLM prompt for Qualification Detection
	// Obtain Specialization
	specialization, err := u.specializationRepo.GetByID(int(specializationResponse.SpecializationTypeID))
	if err != nil {
		log.Printf("Error getting specialization by ID. error: %v", err)
		return err
	}

	// Obtain qualifications
	qualifications, err := u.qualificationRepo.GetAll()
	if err != nil {
		log.Printf("Error getting all job groups. error: %v", err)
		return err
	}
	var qualificationsString string

	for _, qualification := range qualifications {
		qualificationsString += fmt.Sprintf("Qualification name: %s, id: %d; ", qualification.TitleEn, qualification.ID)
	}

	qualificationsString = strings.TrimSuffix(qualificationsString, "; ")

	qualificationDetectionData := prompts_storage.QualificationDetectionData{ResumeText: standarizedResumeText, JobGroupType: jobGroup.TitleEn, SpecializationType: specialization.TitleEn, QualificationTypes: qualificationsString}
	qualificationPrompt, err := pc.GetPrompt(allPrompts.QualificationDetectionPrompt, qualificationDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("QualificationDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the resume
	llmResponse, _, err = u.mistralService.CallMistral(qualificationPrompt, true, mistral.Nemo, "QualificationDetection", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var qualificationResponse QualificationDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &qualificationResponse); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Save the extracted text and rewritten resume in the database
	resume.JobGroupID = &groupDetectionResponse.JobTypeID
	resume.SpecializationID = &specializationResponse.SpecializationTypeID
	resume.QualificationID = &qualificationResponse.QualificationTypeID
	resume.CleanText = formattedText
	resume.RawText = extractedText
	resume.StandarizedText = standarizedResumeText
	resume.FullName = &groupDetectionResponse.FullName
	log.Printf("Saving resume to the database. Resume ID: %v", resume.ID)
	if err := u.resumeRepo.CreateResume(resume); err != nil {
		log.Printf("Error saving resume. error: %v", err)
		return fmt.Errorf("error saving resume: %v", err)
	}

	log.Printf("Starting generating embedding")
	embedding, err := u.mistralService.GenerateEmbedding(resume.StandarizedText)
	if err != nil {
		fmt.Printf("Failed to generate embedding: %v\n", err)
	} else {
		fmt.Println("Generated embedding:", embedding)
	}

	// Convert the []float64 embedding to []float32
	embeddingFloat32 := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingFloat32[i] = float32(v)
	}

	// Save the embedding as a []float32
	resumeEmbedding := &model.ResumeEmbedding{
		ResumeID:  resume.ID,
		Embedding: embeddingFloat32, // Store the []float32 directly
	}

	// Pass both resumeID and embedding to the repository
	if err := u.embeddingRepo.CreateResumeEmbedding(resumeEmbedding.ResumeID, resumeEmbedding.Embedding); err != nil {
		return fmt.Errorf("failed to create resume embedding: %v", err)
	}

	// Log the successful completion of the function
	log.Printf("UploadResume completed successfully for Resume ID: %v", resume.ID)

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
	resume, err := u.resumeRepo.GetByID(resumeID)
	if err != nil {
		return nil, errors.New("resume not found")
	}

	// Construct LLM prompt
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	resumeAnalysisData := prompts_storage.ResumeAnalysisData{ResumeText: resume.CleanText}
	prompt, err := pc.GetPrompt(allPrompts.ResumeAnalysisPrompt, resumeAnalysisData, "", true)
	if err != nil {
		return nil, errors.New("error generating LLM prompt for resume analysis")
	}

	// Call LLM to rewrite the resume in a structured format
	resumeAnalysis, gptCallID, err := u.mistralService.CallMistral(prompt, true, mistral.Nemo, "ResumeAnalysis", 0)
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

// UpdateAllResumeEmbeddings generates embeddings for all resumes without an embedding
func (u *resumeUsecase) UpdateAllResumeEmbeddings() error {
	resumes, err := u.resumeRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to retrieve resumes: %w", err)
	}

	for _, resume := range resumes {
		// Check if embedding exists
		_, err := u.embeddingRepo.GetResumeEmbedding(resume.ID)
		if err == nil {
			log.Printf("Embedding already exists for Resume ID: %v, skipping...", resume.ID)
			continue
		}

		// Generate embedding
		embedding, err := u.mistralService.GenerateEmbedding(resume.StandarizedText)
		if err != nil {
			log.Printf("Error generating embedding for Resume ID: %v, error: %v", resume.ID, err)
			continue
		}

		// Convert embedding from []float64 to []float32
		embeddingFloat32 := make([]float32, len(embedding))
		for i, v := range embedding {
			embeddingFloat32[i] = float32(v)
		}

		// Save embedding
		if err := u.embeddingRepo.CreateResumeEmbedding(resume.ID, embeddingFloat32); err != nil {
			log.Printf("Error saving embedding for Resume ID: %v, error: %v", resume.ID, err)
			continue
		}

		log.Printf("Successfully updated embedding for Resume ID: %v", resume.ID)
	}
	return nil
}
