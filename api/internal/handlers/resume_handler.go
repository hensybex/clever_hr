// internal/handlers/resume_handler.go

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

// Employee uploading a candidate's resume
func (h *ResumeHandler) EmployeeUploadResume(c *gin.Context) {
	// Log that the request has been received
	log.Println("Received request to upload resume from employee")

	// Get the Telegram ID from the form
	telegramIDStr := c.PostForm("tg_id")
	if telegramIDStr == "" {
		log.Println("Error: Telegram ID is missing from the request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telegram ID is required"})
		return
	}
	log.Printf("Processing upload for Telegram ID: %s", telegramIDStr)

	// Retrieve the user by their Telegram ID
	user, err := h.userUsecase.GetUserByTgID(telegramIDStr)
	if err != nil {
		log.Printf("Error: User with Telegram ID %s not found", telegramIDStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	log.Printf("User found: ID = %s", telegramIDStr)
	log.Printf("User found: ID = %d", user.ID)

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
	userID := user.ID
	resume := &model.Resume{
		UploadedByUserID: &userID,
		ResumePDF:        filePath,
	}
	log.Printf("Prepared resume object for user ID %d", userID)

	// Call the usecase to process and save the resume
	if err := h.resumeUsecase.EmployeeUploadResume(resume, filePath); err != nil {
		log.Printf("Error: Failed to process resume for user ID %d", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resume"})
		return
	}
	log.Printf("Resume processed and saved successfully for user ID %d", userID)

	// Success case
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Resume uploaded successfully"})
	log.Println("Resume upload completed successfully")
}

func (h *ResumeHandler) CandidateUploadResume(c *gin.Context) {
	// Retrieve the candidate's Telegram ID from the form data
	telegramIDStr := c.PostForm("tg_id")
	if telegramIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telegram ID is required"})
		return
	}

	// Convert the Telegram ID to a string and retrieve the candidate by their Telegram ID
	user, err := h.userUsecase.GetUserByTgID(telegramIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Retrieve the uploaded resume file
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	// Save the file locally (e.g., in the /tmp directory)
	// Generate a unique file name using UUID and get the file extension
	ext := filepath.Ext(file.Filename)
	uniqueFileName := uuid.New().String() + ext
	log.Printf("Generated unique file name: %s", uniqueFileName)
	filePath := "/app/uploads/" + uniqueFileName
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume file"})
		return
	}

	// Prepare the resume object to be saved
	resume := &model.Resume{
		CandidateID: user.ID, // Use the candidate's ID from the user object
		ResumePDF:   file.Filename,
	}

	// Call the usecase to process and save the resume
	if err := h.resumeUsecase.CandidateUploadResume(resume, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resume"})
		return
	}

	// Success case
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Resume uploaded successfully"})
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
