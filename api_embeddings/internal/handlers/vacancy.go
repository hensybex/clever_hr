package handlers

import (
	"clever_hr_embeddings/internal/models"
	"clever_hr_embeddings/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostVacancyHandler(uc usecases.VacancyUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vacancy models.Vacancy
		if err := c.ShouldBindJSON(&vacancy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := uc.CreateVacancy(&vacancy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Vacancy posted successfully"})
	}
}
