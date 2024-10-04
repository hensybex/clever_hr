// internal/handlers/interview_handler.go

package handlers

import (
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
	if err := c.ShouldBindJSON(&interviewDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	interview := interviewDTO.ToInterviewModel()
	if err := h.interviewUsecase.CreateInterview(&interview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Interview created successfully"})
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

func (h *InterviewHandler) GetInterviewAnalysisResult(c *gin.Context) {
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
