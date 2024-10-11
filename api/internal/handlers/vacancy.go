// internal/handlers/vacancy.go

package handlers

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/usecase"
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

func (h *VacancyHandler) PostVacancy(c *gin.Context) {
	var vacancy model.Vacancy
	if err := c.ShouldBindJSON(&vacancy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.vacancyUsecase.CreateVacancy(&vacancy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vacancy posted successfully"})
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
