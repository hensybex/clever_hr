// handlers/embedding_handler.go

package handlers

import (
	"clever_hr_api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmbeddingHandler struct {
	embeddingUsecase usecase.EmbeddingUsecase
}

func NewEmbeddingHandler(uc usecase.EmbeddingUsecase) *EmbeddingHandler {
	return &EmbeddingHandler{uc}
}

func (h *EmbeddingHandler) GetResumeEmbedding(c *gin.Context) {
	resumeID, err := strconv.Atoi(c.Param("resume_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resume ID"})
		return
	}

	embedding, err := h.embeddingUsecase.GetResumeEmbedding(uint(resumeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"embedding": embedding})
}

func (h *EmbeddingHandler) GetVacancyEmbedding(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("vacancy_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vacancy ID"})
		return
	}

	embedding, err := h.embeddingUsecase.GetVacancyEmbedding(uint(vacancyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"embedding": embedding})
}

func (h *EmbeddingHandler) FindMatchingResumes(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("vacancy_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vacancy ID"})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10 // default value
	}

	embedding, err := h.embeddingUsecase.GetVacancyEmbedding(uint(vacancyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	matches, err := h.embeddingUsecase.FindMatchingResumes(embedding, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}
