// internal/handlers/vacancy.go

package handlers

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VacancyHandler struct {
	vacancyUsecase usecase.VacancyUsecase
}

func NewVacancyHandler(vacancyUsecase usecase.VacancyUsecase) *VacancyHandler {
	return &VacancyHandler{vacancyUsecase}
}

type VacancyDescriptionRequest struct {
	Description string `json:"description" binding:"required"`
}

func (h *VacancyHandler) UploadVacancy(c *gin.Context) {
	// Get the user claims (using the correct key, "id", from JWT middleware)
	user, exists := c.Get("id")
	if !exists {
		log.Println("User not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Extract user object from context
	userObj, ok := user.(*model.User)
	if !ok {
		log.Println("Failed to extract user object from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Read the vacancy description from the client
	var req VacancyDescriptionRequest
	if err := conn.ReadJSON(&req); err != nil {
		log.Println("Read error:", err)
		conn.WriteJSON(dtos.ServerMessage{Status: "Invalid input"})
		return
	}

	// Create a Vacancy object with the description field and user ID
	vacancy := model.Vacancy{
		Description: req.Description,
		UploaderID:  userObj.ID, // Assuming "id" is of type uint
	}

	// Call the usecase to process the vacancy and pass the WebSocket connection
	if err := h.vacancyUsecase.UploadVacancy(&vacancy, conn); err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error: " + err.Error()})
		return
	}

	// Send completion message
	conn.WriteJSON(dtos.ServerMessage{Status: "Vacancy uploaded and resumes matched successfully"})
}

func (h *VacancyHandler) GetVacancies(c *gin.Context) {
	vacancies, err := h.vacancyUsecase.GetAllVacancies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vacancies"})
		return
	}

	c.JSON(http.StatusOK, vacancies)
}

func (h *VacancyHandler) GetVacancyByID(c *gin.Context) {
	vacancyIDStr := c.Param("id")
	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	vacancy, err := h.vacancyUsecase.GetVacancyByID(uint(vacancyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

func (h *VacancyHandler) UpdateVacancyStatus(c *gin.Context) {
	vacancyIDStr := c.Param("id")
	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.vacancyUsecase.UpdateVacancyStatus(uint(vacancyID), statusUpdate.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vacancy status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vacancy status updated successfully"})
}

// UpdateAllVacancyEmbeddings updates embeddings for all vacancies
func (h *VacancyHandler) UpdateAllVacancyEmbeddings(c *gin.Context) {
	err := h.vacancyUsecase.UpdateAllVacancyEmbeddings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All vacancy embeddings updated successfully"})
}
