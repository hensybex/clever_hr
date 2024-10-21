// usecase/match.go

package usecase

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/mistral"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/prompts"
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/service"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type VacancyResumeMatchUsecase interface {
	MatchVacancyWithResumes(vacancyID uint, conn *websocket.Conn) error
	MatchVacancyWithResume(vacancyID, resumeID uint) error
	GetMatchesByVacancyID(vacancyID uint) ([]model.VacancyResumeMatch, error)
	GetMatchByID(matchID uint) (*model.VacancyResumeMatch, error)
}

type MatchUsecase struct {
	embeddingRepo          repository.EmbeddingRepository
	vacancyResumeMatchRepo repository.VacancyResumeMatchRepository
	vacancyRepo            repository.VacancyRepository
	resumeRepo             repository.ResumeRepository
	mistralService         mistral.MistralService
}

func NewMatchUsecase(
	embeddingRepo repository.EmbeddingRepository,
	vacancyResumeMatchRepo repository.VacancyResumeMatchRepository,
	vacancyRepo repository.VacancyRepository,
	resumeRepo repository.ResumeRepository,
	mistralService mistral.MistralService,
) VacancyResumeMatchUsecase {
	return &MatchUsecase{
		embeddingRepo:          embeddingRepo,
		vacancyResumeMatchRepo: vacancyResumeMatchRepo,
		vacancyRepo:            vacancyRepo,
		resumeRepo:             resumeRepo,
		mistralService:         mistralService,
	}
}

// MatchVacancyWithResumes matches a vacancy with resumes using vector similarity
func (u *MatchUsecase) MatchVacancyWithResumes(vacancyID uint, conn *websocket.Conn) error {
	// Get the vacancy embedding
	vacancyEmbeddingRecord, err := u.embeddingRepo.GetVacancyEmbedding(vacancyID)
	if err != nil {
		return fmt.Errorf("failed to get vacancy embedding: %w", err)
	}

	// Send status update
	conn.WriteJSON(dtos.ServerMessage{Status: "Fetching matching resumes"})

	// Find matching resumes using the embedding directly
	matches, err := u.embeddingRepo.FindMatchingResumes(vacancyEmbeddingRecord, 10)
	if err != nil {
		return fmt.Errorf("failed to find matching resumes: %w", err)
	}

	// Get vacancy details
	vacancy, err := u.vacancyRepo.GetByID(vacancyID)
	if err != nil {
		return fmt.Errorf("failed to get vacancy details: %w", err)
	}

	// Initialize prompts
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()

	// Iterate over the top matches
	for idx, match := range matches {
		conn.WriteJSON(dtos.ServerMessage{Status: fmt.Sprintf("Analyzing resume %d of %d", idx+1, len(matches))})

		// Get resume details
		resume, err := u.resumeRepo.GetByID(match.ResumeID)
		if err != nil {
			log.Printf("Error fetching resume ID %d: %v", match.ResumeID, err)
			continue // Skip this resume and proceed to the next
		}

		// Prepare data for the prompt
		resumeVacancyData := prompts_storage.ResumeVacancyAnalysisData{
			ResumeText:   resume.CleanText,
			VacancyText:  vacancy.Description,
			VacancyTitle: vacancy.Title,
			// Include other relevant vacancy fields if needed
		}

		// Construct the prompt
		prompt, err := pc.GetPrompt(allPrompts.ResumeVacancyAnalysisPrompt, resumeVacancyData, "", true)
		if err != nil {
			log.Printf("Error generating LLM prompt: %v", err)
			continue // Skip this resume and proceed to the next
		}

		// Call the LLM service
		llmResponse, gptCallID, err := u.mistralService.CallMistral(prompt, true, mistral.Largest, "ResumeVacancyAnalysis", 0)
		if err != nil {
			log.Printf("Error calling LLM service: %v", err)
			continue // Skip this resume and proceed to the next
		}

		// Parse LLM response into structured data
		var analysisResponse struct {
			RelevantWorkExperience               model.AnalysisField `json:"RelevantWorkExperience"`
			TechnicalSkillsAndProficiencies      model.AnalysisField `json:"TechnicalSkillsAndProficiencies"`
			EducationAndCertifications           model.AnalysisField `json:"EducationAndCertifications"`
			SoftSkillsAndCulturalFit             model.AnalysisField `json:"SoftSkillsAndCulturalFit"`
			LanguageAndCommunicationSkills       model.AnalysisField `json:"LanguageAndCommunicationSkills"`
			ProblemSolvingAndAnalyticalAbilities model.AnalysisField `json:"ProblemSolvingAndAnalyticalAbilities"`
			AdaptabilityAndLearningCapacity      model.AnalysisField `json:"AdaptabilityAndLearningCapacity"`
			LeadershipAndManagementExperience    model.AnalysisField `json:"LeadershipAndManagementExperience"`
			MotivationAndCareerObjectives        model.AnalysisField `json:"MotivationAndCareerObjectives"`
			AdditionalQualificationsAndValueAdds model.AnalysisField `json:"AdditionalQualificationsAndValueAdds"`
		}

		if err := json.Unmarshal([]byte(llmResponse), &analysisResponse); err != nil {
			log.Printf("Error parsing LLM response: %v", err)
			continue // Skip this resume and proceed to the next
		}

		// Create or update the VacancyResumeMatch record
		vacancyResumeMatch := &model.VacancyResumeMatch{
			VacancyID: vacancyID,
			ResumeID:  match.ResumeID,
			Score:     match.Score,

			// Analysis fields
			RelevantWorkExperience:               analysisResponse.RelevantWorkExperience,
			TechnicalSkillsAndProficiencies:      analysisResponse.TechnicalSkillsAndProficiencies,
			EducationAndCertifications:           analysisResponse.EducationAndCertifications,
			SoftSkillsAndCulturalFit:             analysisResponse.SoftSkillsAndCulturalFit,
			LanguageAndCommunicationSkills:       analysisResponse.LanguageAndCommunicationSkills,
			ProblemSolvingAndAnalyticalAbilities: analysisResponse.ProblemSolvingAndAnalyticalAbilities,
			AdaptabilityAndLearningCapacity:      analysisResponse.AdaptabilityAndLearningCapacity,
			LeadershipAndManagementExperience:    analysisResponse.LeadershipAndManagementExperience,
			MotivationAndCareerObjectives:        analysisResponse.MotivationAndCareerObjectives,
			AdditionalQualificationsAndValueAdds: analysisResponse.AdditionalQualificationsAndValueAdds,
			AnalysisStatus:                       "finished",
			GPTCallID:                            &gptCallID,
		}

		if err := u.vacancyResumeMatchRepo.CreateOrUpdateMatch(vacancyResumeMatch); err != nil {
			log.Printf("Error saving vacancy-resume match: %v", err)
			continue // Skip this resume and proceed to the next
		}
	}

	// Send completion status
	conn.WriteJSON(dtos.ServerMessage{Status: "Analysis completed"})

	return nil
}

// MatchVacancyWithResume matches a specific vacancy with a specific resume
func (u *MatchUsecase) MatchVacancyWithResume(vacancyID, resumeID uint) error {
	// Get the vacancy embedding
	vacancyEmbeddingRecord, err := u.embeddingRepo.GetVacancyEmbedding(vacancyID)
	if err != nil {
		return fmt.Errorf("failed to get vacancy embedding: %w", err)
	}

	// Get the resume embedding
	resumeEmbeddingRecord, err := u.embeddingRepo.GetResumeEmbedding(resumeID)
	if err != nil {
		return fmt.Errorf("failed to get resume embedding: %w", err)
	}

	// Calculate similarity using cosine similarity in Go
	score, err := service.CalculateCosineSimilarity(
		vacancyEmbeddingRecord, // Pass the float32 slice directly
		resumeEmbeddingRecord,  // Pass the float32 slice directly
	)
	if err != nil {
		return fmt.Errorf("failed to calculate similarity: %w", err)
	}

	// Insert the match into the database
	vacancyResumeMatch := &model.VacancyResumeMatch{
		VacancyID: vacancyID,
		ResumeID:  resumeID,
		Score:     score,
	}
	if err := u.vacancyResumeMatchRepo.CreateMatch(vacancyResumeMatch); err != nil {
		return fmt.Errorf("failed to create vacancy-resume match: %w", err)
	}

	return nil
}

func (u *MatchUsecase) GetMatchesByVacancyID(vacancyID uint) ([]model.VacancyResumeMatch, error) {
	return u.vacancyResumeMatchRepo.GetMatchesByVacancyID(vacancyID)
}

func (u *MatchUsecase) GetMatchByID(matchID uint) (*model.VacancyResumeMatch, error) {
	return u.vacancyResumeMatchRepo.GetMatchByID(matchID)
}
