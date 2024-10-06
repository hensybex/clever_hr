// internal/handlers/interview_handler.go

package handlers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type InterviewHandler struct {
	interviewUsecase usecase.InterviewUsecase
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewInterviewHandler(interviewUsecase usecase.InterviewUsecase) *InterviewHandler {
	return &InterviewHandler{interviewUsecase}
}

func (h *InterviewHandler) CreateInterview(c *gin.Context) {
	var interviewDTO dtos.CreateInterviewDTO

	// Log the incoming request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Printf("Request Body: %s", string(bodyBytes))

	// Reset the body for ShouldBindJSON to read it again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Bind the JSON input to interviewDTO (without ResumeID for now)
	if err := c.ShouldBindJSON(&interviewDTO); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the usecase and get the created interview ID
	interviewID, err := h.interviewUsecase.CreateInterview(interviewDTO)
	if err != nil {
		if err == errors.New("no resume found") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Resume not found",
				"message": "Пожалуйста, загрузите резюме, прежде чем начинать собеседование.",
			})
		} else {
			log.Printf("Error creating interview: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while creating the interview"})
		}
		return
	}

	// Success response with the interview ID
	log.Printf("Interview created successfully with ID: %d", interviewID)
	c.JSON(http.StatusCreated, gin.H{
		"success":      true,
		"message":      "Interview created successfully",
		"interview_id": interviewID, // Add the interview ID to the response
	})
}

func (h *InterviewHandler) AnalyseInterviewMessageWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set websocket upgrade"})
		return
	}
	defer conn.Close()

	for {
		var clientMsg dtos.ClientMessage
		if err := conn.ReadJSON(&clientMsg); err != nil {
			log.Println("Read error:", err)
			break
		}

		go h.interviewUsecase.AnalyseInterviewMessageWebsocket(clientMsg, conn)
	}
}

func (h *InterviewHandler) RunFullInterviewAnalysis(c *gin.Context) {
	interviewIDStr := c.Param("interview_id")
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interview ID"})
		return
	}

	result, err := h.interviewUsecase.RunFullInterviewAnalysis(uint(interviewID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})

}

func (h *InterviewHandler) GetInterviewAnalysisResultByInterviewID(c *gin.Context) {
	interviewIDStr := c.Param("interview_id")
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interview ID"})
		return
	}

	result, err := h.interviewUsecase.GetInterviewAnalysisResult(uint(interviewID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview analysis result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
}

func (h *InterviewHandler) GetInterviewAnalysisResult(c *gin.Context) {
	interviewIDStr := c.Param("interview_id")
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interview ID"})
		return
	}

	result, err := h.interviewUsecase.GetOneByID(uint(interviewID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview analysis result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
}