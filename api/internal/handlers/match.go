// handlers/match.go

package handlers

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MatchHandler struct {
	matchUsecase usecase.VacancyResumeMatchUsecase
}

func NewMatchHandler(matchUsecase usecase.VacancyResumeMatchUsecase) *MatchHandler {
	return &MatchHandler{matchUsecase}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Adjust the origin check as per your requirements
		return true
	},
}

func (h *MatchHandler) MatchVacancyWithResumes(c *gin.Context) {
	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	vacancyIDStr := c.Param("vacancy_id")
	vacancyID, err := strconv.ParseUint(vacancyIDStr, 10, 64)
	if err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Invalid vacancy_id"})
		return
	}

	// Start the matching process
	err = h.matchUsecase.MatchVacancyWithResumes(uint(vacancyID), conn)
	if err != nil {
		conn.WriteJSON(dtos.ServerMessage{Status: "Error: " + err.Error()})
		return
	}
}

// MatchVacancyWithResume matches a specific vacancy with a specific resume
func (h *MatchHandler) MatchVacancyWithResume(c *gin.Context) {
	vacancyIDStr := c.Param("vacancy_id")
	vacancyID, err := strconv.ParseUint(vacancyIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vacancy_id"})
		return
	}

	resumeIDStr := c.Param("resume_id")
	resumeID, err := strconv.ParseUint(resumeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resume_id"})
		return
	}

	if err := h.matchUsecase.MatchVacancyWithResume(uint(vacancyID), uint(resumeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *MatchHandler) GetVacancyResumeMatches(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("vacancy_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	matches, err := h.matchUsecase.GetMatchesByVacancyID(uint(vacancyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}

func (h *MatchHandler) GetVacancyResumeMatchByID(c *gin.Context) {
	matchID := c.Param("match_id")

	// Convert matchID to an integer
	id, err := strconv.Atoi(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Fetch the VacancyResumeMatch from the database
	match, err := h.matchUsecase.GetMatchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch match"})
		return
	}

	c.JSON(http.StatusOK, match)
}
