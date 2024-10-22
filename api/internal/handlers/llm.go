// internal/handlers/llm_handler.go

package handlers

import (
	"clever_hr_api/internal/mistral"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/prompts"
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/tools"
	"clever_hr_api/internal/usecase"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LLMHandler struct {
	messageUsecase         usecase.MessageUsecase
	vacancyResumeMatchRepo repository.VacancyResumeMatchRepository
	mistralService         *mistral.MistralService
}

func NewLLMHandler(messageUsecase usecase.MessageUsecase, vacancyResumeMatchRepo repository.VacancyResumeMatchRepository, mistralService *mistral.MistralService) *LLMHandler {
	return &LLMHandler{
		messageUsecase:         messageUsecase,
		vacancyResumeMatchRepo: vacancyResumeMatchRepo,
		mistralService:         mistralService,
	}
}

type InterviewResponse struct {
	Response string `json:"Response"`
}

type StartingResponse struct {
	StartMessage string `json:"StartMessage"`
}

// GenerateInterviewAnswer handles POST /llm/generate_interview_answer
func (h *LLMHandler) GenerateInterviewAnswer(c *gin.Context) {
	var request struct {
		TgLogin     string `json:"tg_login" binding:"required"`
		LastMessage string `json:"last_message" binding:"required"` // New field for last human message
	}

	// Parse incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Fetch messages by tg_login
	messages, err := h.messageUsecase.GetMessagesByTgLogin(request.TgLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	// Format previous messages
	formattedMessages := tools.FormatMessages(messages)

	// Append the last message from the user
	formattedMessages += "Human reply: " + request.LastMessage + "\n"

	// Save the last human message to the database
	humanMessage := model.Message{
		TgLogin:    request.TgLogin,
		Message:    request.LastMessage,
		SentByUser: true, // Mark as a human-sent message
		CreatedAt:  time.Now(),
	}

	if err := h.messageUsecase.CreateMessage(&humanMessage); err != nil {
		log.Printf("Error saving human message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save human message"})
		return
	}
	log.Printf("Human message saved successfully.")

	// Create prompts for the LLM
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	interviewData := prompts_storage.InterviewMessageData{PreviousDialog: formattedMessages}
	prompt, err := pc.GetPrompt(allPrompts.InterviewMessagePrompt, interviewData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate prompt for LLM"})
		return
	}
	log.Printf("Interview prompt constructed successfully.")

	// Call the LLM service with the generated prompt
	llmResponse, _, err := h.mistralService.CallMistral(prompt, false, mistral.Largest, "InterviewMessage", 0)
	if err != nil {
		log.Printf("Error calling LLM service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call LLM service"})
		return
	}
	log.Printf("LLM response received successfully.")

	// Parse LLM response into a structured interview response
	var interviewResponse InterviewResponse
	if err := json.Unmarshal([]byte(llmResponse), &interviewResponse); err != nil {
		log.Printf("Error parsing LLM response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse LLM response"})
		return
	}

	// Create a new message with the AI response and save it to the database
	aiMessage := model.Message{
		TgLogin:    request.TgLogin,
		Message:    interviewResponse.Response,
		SentByUser: false, // AI-generated message
		CreatedAt:  time.Now(),
	}

	if err := h.messageUsecase.CreateMessage(&aiMessage); err != nil {
		log.Printf("Error saving AI message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save AI message"})
		return
	}
	log.Printf("AI message saved successfully.")

	// Return the generated AI message to the client
	c.JSON(http.StatusOK, gin.H{"message": interviewResponse.Response})
}

// GenerateFirstMessage handles POST /llm/generate_first_message
func (h *LLMHandler) GenerateFirstMessage(c *gin.Context) {
	var request struct {
		MatchID string `json:"match_id" binding:"required"`
		TgLogin string `json:"tg_login" binding:"required"`
	}

	// Parse the incoming request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convert MatchID from string to uint
	matchID, err := strconv.ParseUint(request.MatchID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Fetch the VacancyResumeMatch by ID from the usecase
	vacancyResumeMatch, err := h.vacancyResumeMatchRepo.GetMatchByID(uint(matchID))
	if err != nil {
		log.Printf("Error fetching VacancyResumeMatch by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch match data"})
		return
	}

	// Extract the analysis fields from the VacancyResumeMatch object
	analysisData := formatAnalysisFields(vacancyResumeMatch)

	// Create prompts for the LLM using the analysis data
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	promptData := prompts_storage.FirstMessageData{AnalysisData: analysisData}
	prompt, err := pc.GetPrompt(allPrompts.FirstMessagePrompt, promptData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate prompt for LLM"})
		return
	}
	log.Printf("First message prompt constructed successfully.")

	// Call the LLM service with the generated prompt
	llmResponse, _, err := h.mistralService.CallMistral(prompt, false, mistral.Largest, "FirstMessage", 0)
	if err != nil {
		log.Printf("Error calling LLM service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call LLM service"})
		return
	}
	log.Printf("LLM response received successfully.")

	// Parse LLM response into a structured interview response
	var startingResponse StartingResponse
	if err := json.Unmarshal([]byte(llmResponse), &startingResponse); err != nil {
		log.Printf("Error parsing LLM response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse LLM response"})
		return
	}

	// Create a new message with the AI response and save it to the database
	aiMessage := model.Message{
		TgLogin:    request.TgLogin,
		Message:    startingResponse.StartMessage,
		SentByUser: false, // AI-generated message
		CreatedAt:  time.Now(),
	}

	if err := h.messageUsecase.CreateMessage(&aiMessage); err != nil {
		log.Printf("Error saving AI message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save AI message"})
		return
	}

	// Return the generated message to the client
	c.JSON(http.StatusOK, gin.H{"message": startingResponse.StartMessage})
}

// Helper function to format analysis fields from VacancyResumeMatch
func formatAnalysisFields(match *model.VacancyResumeMatch) string {
	var builder strings.Builder
	builder.WriteString("Relevant Work Experience: " + match.RelevantWorkExperience.Overview + "\n")
	builder.WriteString("Technical Skills and Proficiencies: " + match.TechnicalSkillsAndProficiencies.Overview + "\n")
	builder.WriteString("Education and Certifications: " + match.EducationAndCertifications.Overview + "\n")
	builder.WriteString("Soft Skills and Cultural Fit: " + match.SoftSkillsAndCulturalFit.Overview + "\n")
	builder.WriteString("Language and Communication Skills: " + match.LanguageAndCommunicationSkills.Overview + "\n")
	builder.WriteString("Problem-Solving and Analytical Abilities: " + match.ProblemSolvingAndAnalyticalAbilities.Overview + "\n")
	builder.WriteString("Adaptability and Learning Capacity: " + match.AdaptabilityAndLearningCapacity.Overview + "\n")
	builder.WriteString("Leadership and Management Experience: " + match.LeadershipAndManagementExperience.Overview + "\n")
	builder.WriteString("Motivation and Career Objectives: " + match.MotivationAndCareerObjectives.Overview + "\n")
	builder.WriteString("Additional Qualifications and Value Adds: " + match.AdditionalQualificationsAndValueAdds.Overview + "\n")

	return builder.String()
}
