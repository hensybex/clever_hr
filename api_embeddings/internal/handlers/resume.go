package handlers

import (
	"clever_hr_embeddings/internal/models"
	"clever_hr_embeddings/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadResumeHandler(uc usecases.ResumeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resume models.Resume
		if err := c.ShouldBindJSON(&resume); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := uc.CreateResume(&resume); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded successfully"})
	}
}
