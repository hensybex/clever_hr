package handlers

import (
	"clever_hr_embeddings/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetBestCandidatesHandler(uc usecases.CandidateMatchingUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		vacancyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
			return
		}

		candidates, err := uc.GetBestCandidates(uint(vacancyID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, candidates)
	}
}
