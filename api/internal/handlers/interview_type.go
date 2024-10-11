// internal/handlers/interview_type.go

package handlers

import (
	"net/http"

	"clever_hr_api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type InterviewTypeHandler struct {
	usecase usecase.InterviewTypeUsecase
}

func NewInterviewTypeHandler(usecase usecase.InterviewTypeUsecase) *InterviewTypeHandler {
	return &InterviewTypeHandler{usecase}
}

// ListInterviewTypes handles the request to list all interview types
func (h *InterviewTypeHandler) ListInterviewTypes(c *gin.Context) {
	interviewTypes, err := h.usecase.ListInterviewTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve interview types"})
		return
	}

	c.JSON(http.StatusOK, interviewTypes)
}
