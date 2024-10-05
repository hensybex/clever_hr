// internal/handlers/user_handler.go

package handlers

import (
	"log"
	"net/http"
	"strconv"

	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase      usecase.UserUsecase
	candidateUsecase usecase.CandidateUsecase
}

func NewUserHandler(
	userUsecase usecase.UserUsecase,
	candidateUsecase usecase.CandidateUsecase,
) *UserHandler {
	return &UserHandler{
		userUsecase,
		candidateUsecase,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var userDTO dtos.CreateUserDTO

	// Log the incoming request
	log.Println("INFO: Received request to create a new user.")

	// Bind the incoming JSON to userDTO and check for errors
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		log.Printf("ERROR: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("INFO: Request payload: %+v", userDTO)

	// Convert DTO to user model
	user := userDTO.ToUserModel()

	// Log the user model data being created (consider redacting sensitive information if necessary)
	log.Printf("INFO: Creating user with details: %+v", user)

	// Attempt to create the user via usecase
	err := h.userUsecase.CreateUser(&user)
	if err != nil {
		log.Printf("ERROR: Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log success and return the response
	log.Printf("INFO: User created successfully: %+v", user)
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User created successfully"})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SwitchUserType(c *gin.Context) {
	telegramUserIDStr := c.Param("user_id")
	telegramUserID, err := strconv.Atoi(telegramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	telegramID := strconv.FormatUint(uint64(telegramUserID), 10)
	user, err := h.userUsecase.GetUserByTgID(telegramID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.userUsecase.SwitchUserType(uint(user.ID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch user type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Тип пользователя успешно изменен"})
}

func (h *UserHandler) GetCandidatesByUserID(c *gin.Context) {
	// Extract user_id from URL parameters
	userIDStr := c.Param("user_id")

	// Call the usecase to get the list of candidates uploaded by this user
	candidates, err := h.userUsecase.GetCandidatesByTGID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve candidates"})
		return
	}

	// Return the list of candidates
	c.JSON(http.StatusOK, gin.H{"candidates": candidates})
}

func (h *UserHandler) GetUserRoleByTgID(c *gin.Context) {
	// Get tg_id from the URL parameter
	tgID := c.Param("user_id")

	// Call the use case to get the user role
	role, err := h.userUsecase.GetUserRoleByTgID(tgID)
	if err != nil {
		log.Printf("Error fetching user role for tg_id %s: %v", tgID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user role"})
		return
	}

	// If role is not found, return a 404 response
	if role == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the role as a JSON response
	c.JSON(http.StatusOK, gin.H{"role": role})
}
