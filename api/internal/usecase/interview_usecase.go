// internal/usecase/interview_usecase.go

package usecase

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/prompts"
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/service"
	"errors"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type InterviewUsecase interface {
	CreateInterview(interviewDTO dtos.CreateInterviewDTO) (uint, error)
	AnalyseInterviewMessageWebsocket(clientMsg dtos.ClientMessage, conn *websocket.Conn)
	RunFullInterviewAnalysis(interviewID uint) (*model.InterviewAnalysisResult, error)
	GetInterviewAnalysisResult(interviewID uint) (*model.InterviewAnalysisResult, error)
}

type interviewUsecase struct {
	interviewRepo     repository.InterviewRepository
	interviewTypeRepo repository.InterviewTypeRepository
	messageRepo       repository.InterviewMessageRepository
	analysisRepo      repository.InterviewAnalysisResultRepository
	resumeRepo        repository.ResumeRepository
	userRepo          repository.UserRepository
	candidateRepo     repository.CandidateRepository
	mistralService    service.MistralService
}

func NewInterviewUsecase(
	interviewRepo repository.InterviewRepository,
	interviewTypeRepo repository.InterviewTypeRepository,
	messageRepo repository.InterviewMessageRepository,
	analysisRepo repository.InterviewAnalysisResultRepository,
	resumeRepo repository.ResumeRepository,
	userRepo repository.UserRepository,
	candidateRepo repository.CandidateRepository,
	mistralService service.MistralService,
) InterviewUsecase {
	return &interviewUsecase{
		interviewRepo,
		interviewTypeRepo,
		messageRepo,
		analysisRepo,
		resumeRepo,
		userRepo,
		candidateRepo,
		mistralService,
	}
}

func (u *interviewUsecase) CreateInterview(interviewDTO dtos.CreateInterviewDTO) (uint, error) {
	// Check if the user has uploaded a resume based on tg_id
	user, err := u.userRepo.GetUserByTgID(strconv.FormatUint(uint64(interviewDTO.TgID), 10))
	if err != nil {
		return 0, err
	}
	candidate, err := u.candidateRepo.GetCandidateByUploadedID(user.ID)
	if err != nil {
		return 0, err
	}
	resume, err := u.resumeRepo.GetResumeByCandidateID(candidate.ID)
	if err != nil {
		return 0, err
	}
	if resume == nil {
		// No resume found, return a custom error
		return 0, errors.New("no resume found")
	}

	// Convert the DTO to the Interview model
	interview := model.Interview{ResumeID: resume.ID, InterviewTypeID: interviewDTO.InterviewTypeID, Status: model.InterviewStatus(model.InProgress)}

	// Create the interview using the repository
	if err := u.interviewRepo.CreateInterview(&interview); err != nil {
		return 0, err
	}

	return interview.ID, nil
}

func (u *interviewUsecase) AnalyseInterviewMessageWebsocket(clientMsg dtos.ClientMessage, conn *websocket.Conn) {
	// Send initial status
	conn.WriteJSON(dtos.ServerMessage{Status: "Processing started"})

	// Store the new message
	newMessage := &model.InterviewMessage{
		InterviewID: uint(clientMsg.InterviewID),
		MessageText: clientMsg.Message,
		SentBy:      model.SentByCandidate, // assuming candidate, adapt as needed
	}
	if err := u.messageRepo.CreateMessage(newMessage); err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error creating message"})
		return
	}

	// Fetch all previous messages for the interview
	previousMessages, err := u.messageRepo.GetMessagesByInterviewID(uint(clientMsg.InterviewID))
	if err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error fetching previous messages"})
		return
	}

	// Check if the number of previous messages exceeds 6
	if len(previousMessages) > 6 {
		// End the interview if message count exceeds 6
		conn.WriteJSON(dtos.ServerMessage{Status: "End of interview"})
		// Trigger interview analysis
		if _, err := u.RunFullInterviewAnalysis(uint(clientMsg.InterviewID)); err != nil {
			log.Printf("Error triggering interview analysis: %v", err)
		}
		return
	}

	// Fetch the interview details
	interview, err := u.interviewRepo.GetInterviewByID(uint(clientMsg.InterviewID))
	if err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error getting interview"})
		return
	}

	interviewType, err := u.interviewTypeRepo.GetInterviewTypeByID(interview.InterviewTypeID)
	if err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error getting interview type"})
		return
	}

	// Convert previous messages into formatted string
	formattedMessages := service.FormatMessages(previousMessages)

	// Construct prompt for the LLM using the previous messages and the current message
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	dialogData := prompts_storage.InterviewMessageAnalysisData{PreviousMessages: formattedMessages, LatestMessage: clientMsg.Message, InterviewType: interviewType.Name}
	prompt, err := pc.GetPrompt(allPrompts.InterviewMessageAnalysisPrompt, dialogData, "Russian", true)
	if err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error generating prompt"})
		return
	}

	// Call LLM
	fullResponse, gptCallID, err := u.mistralService.CallMistralStream(prompt, true, service.Nemo, "DialogAnalysis", 0, conn)
	if err != nil {
		log.Printf("Error calling LLM for summary: %v", err)
		conn.WriteJSON(dtos.ServerMessage{Status: "Error generating interview response"})
		return
	}

	// Store the new message
	receivedResponseMessage := &model.InterviewMessage{
		InterviewID: uint(clientMsg.InterviewID),
		MessageText: fullResponse,
		SentBy:      model.SentByLLM,
		GPTCallID:   &gptCallID,
	}
	if err := u.messageRepo.CreateMessage(receivedResponseMessage); err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error creating message"})
		return
	}

	conn.WriteJSON(dtos.ServerMessage{Status: "Interview message completed"})
}

func (u *interviewUsecase) RunFullInterviewAnalysis(interviewID uint) (*model.InterviewAnalysisResult, error) {
	// Fetch all interview messages
	messages, err := u.messageRepo.GetMessagesByInterviewID(interviewID)
	if err != nil {
		return nil, errors.New("failed to fetch interview messages")
	}

	interview, err := u.interviewRepo.GetInterviewByID(interviewID)
	if err != nil {
		return nil, errors.New("failed to fetch interview")
	}

	interviewType, err := u.interviewTypeRepo.GetInterviewTypeByID(interview.InterviewTypeID)
	if err != nil {
		return nil, errors.New("failed to fetch interview type")
	}

	// Format messages into a string
	formattedMessages := service.FormatMessages(messages)

	// Construct prompt for the LLM using the previous messages and the current message
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	interviewAnalysisData := prompts_storage.FullInterviewAnalysisData{InterviewMessages: formattedMessages, InterviewType: interviewType.Name}
	prompt, err := pc.GetPrompt(allPrompts.FullInterviewAnalysisPrompt, interviewAnalysisData, "", true)
	if err != nil {
		return nil, errors.New("failed to generate prompt")
	}

	// Call LLM
	fullResponse, gptCallID, err := u.mistralService.CallMistral(prompt, true, service.Nemo, "DialogAnalysis", 0)
	if err != nil {
		log.Printf("Error calling LLM for summary: %v", err)
		return nil, errors.New("llm failure")
	}

	// Simulating full interview analysis result
	analysisResult := &model.InterviewAnalysisResult{
		InterviewID:  interviewID,
		GPTCallID:    &gptCallID,
		ResultStatus: model.ResultFinished,
		Result:       fullResponse,
	}

	// Save the analysis result
	if err := u.analysisRepo.CreateInterviewAnalysisResult(analysisResult); err != nil {
		return nil, err
	}

	return analysisResult, nil
}

func (u *interviewUsecase) GetInterviewAnalysisResult(interviewID uint) (*model.InterviewAnalysisResult, error) {
	return u.analysisRepo.GetInterviewAnalysisResultByInterviewID(interviewID)
}
