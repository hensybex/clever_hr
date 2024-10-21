// internal/handlers/resume.go

package handlers

import (
	//"log"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"clever_hr_api/internal/model"
	"clever_hr_api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResumeHandler struct {
	resumeUsecase               usecase.ResumeUsecase
	userUsecase                 usecase.UserUsecase
	resumeAnalysisResultUsecase usecase.ResumeAnalysisResultUsecase
}

func NewResumeHandler(
	resumeUsecase usecase.ResumeUsecase,
	userUsecase usecase.UserUsecase,
	resumeAnalysisResultUsecase usecase.ResumeAnalysisResultUsecase,
) *ResumeHandler {
	return &ResumeHandler{
		resumeUsecase,
		userUsecase,
		resumeAnalysisResultUsecase,
	}
}

func (h *ResumeHandler) GetResumeByID(c *gin.Context) {
	idStr := c.Param("resume_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resume ID"})
		return
	}

	// Call the usecase to fetch the resume by ID
	resume, err := h.resumeUsecase.GetResumeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	// Return the resume details in the response
	c.JSON(http.StatusOK, gin.H{"resume": resume})
}

// Employee uploading a candidate's resume
func (h *ResumeHandler) UploadResume(c *gin.Context) {
	// Log that the request has been received
	log.Println("Received request to upload resume from employee")

	// Get the Telegram ID from the form
	uploadedFrom := c.PostForm("uploaded_from")
	if uploadedFrom == "" {
		log.Println("Error: UploadedFrom is missing from the request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded From is required"})
		return
	}

	// Retrieve the uploaded resume file
	file, err := c.FormFile("resume")
	if err != nil {
		log.Println("Error: Resume file is missing from the request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}
	log.Printf("Received file: %s", file.Filename)

	// Generate a unique file name using UUID and get the file extension
	ext := filepath.Ext(file.Filename)
	uniqueFileName := uuid.New().String() + ext
	log.Printf("Generated unique file name: %s", uniqueFileName)

	// Save the file locally (e.g., in the /tmp directory)
	filePath := "/app/uploads/" + uniqueFileName
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error: Failed to save file %s at path %s", file.Filename, filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume file"})
		return
	}
	log.Printf("File saved successfully at: %s", filePath)

	// Prepare the resume object to be saved
	resume := model.Resume{
		UploadedFrom: uploadedFrom,
		PDFPath:      filePath,
	}

	// Call the usecase to process and save the resume
	if err := h.resumeUsecase.UploadResume(&resume, filePath); err != nil {
		log.Printf("Error: Failed to upload resume")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resume"})
		return
	}
	log.Printf("Resume processed and saved successfully")

	// Success case
	c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded successfully"})
	log.Println("Resume upload completed successfully")
}

func (h *ResumeHandler) RunResumeAnalysis(c *gin.Context) {
	idStr := c.Param("resume_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resume ID"})
		return
	}

	result, err := h.resumeUsecase.RunResumeAnalysis(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
}

func (h *ResumeHandler) GetResumeAnalysisResult(c *gin.Context) {
	idStr := c.Param("resume_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resume ID"})
		return
	}

	result, err := h.resumeAnalysisResultUsecase.GetResumeAnalysisResultByResumeID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume analysis result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
}

// UpdateAllResumeEmbeddings updates embeddings for all resumes
func (h *ResumeHandler) UpdateAllResumeEmbeddings(c *gin.Context) {
	err := h.resumeUsecase.UpdateAllResumeEmbeddings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All resume embeddings updated successfully"})
}
