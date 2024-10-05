// internal/handlers/candidate_handler.go

package handlers

import (
	"clever_hr_api/internal/usecase"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CandidateHandler struct {
	candidateUsecase usecase.CandidateUsecase
}

func NewCandidateHandler(candidateUsecase usecase.CandidateUsecase) *CandidateHandler {
	return &CandidateHandler{candidateUsecase}
}

func (h *CandidateHandler) GetCandidateInfo(c *gin.Context) {
	candidateIDStr := c.Param("candidate_id")
	candidateID, err := strconv.Atoi(candidateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}

	// Fetch candidate details (without resume)
	candidateInfo, err := h.candidateUsecase.GetCandidateInfo(uint(candidateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidate info"})
		return
	}

	// Return the candidate metadata as JSON (without the resume PDF)
	c.JSON(http.StatusOK, candidateInfo)
}

func (h *CandidateHandler) GetCandidateInfoByTgID(c *gin.Context) {
	tgIDStr := c.Param("candidate_id")

	candidateInfo, err := h.candidateUsecase.GetCandidateInfoByTgID(tgIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidate info"})
		return
	}

	// Return the candidate metadata as JSON (without the resume PDF)
	c.JSON(http.StatusOK, candidateInfo)
}

func (h *CandidateHandler) GetResume(c *gin.Context) {
	// Get the candidate_id and resume_path from the query parameters
	candidateID := c.Param("candidate_id")
	resumePath := c.Query("resume_path") // Retrieve resume_path from query parameter

	if resumePath == "" {
		log.Printf("Missing resume_path for candidate %s", candidateID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing resume_path"})
		return
	}

	// Log the path to ensure it's correct
	log.Printf("Attempting to open file at path: %s", resumePath)

	// Open the PDF file
	file, err := os.Open(resumePath)
	if err != nil {
		log.Printf("Error opening file at path %s: %v", resumePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open resume PDF"})
		return
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info for %s: %v", resumePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file info"})
		return
	}

	log.Printf("Serving file: %s with size: %d bytes", fileInfo.Name(), fileInfo.Size())

	// Set headers for returning the file as a download
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Stream the file back to the client
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
}
