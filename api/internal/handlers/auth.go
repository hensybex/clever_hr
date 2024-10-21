// internal/handlers/auth.go

package handlers

import (
	"clever_hr_api/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginVals struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&loginVals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Authenticate user
	user, err := h.authService.Authenticate(loginVals.Username, loginVals.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Generate JWT token
	token, expire, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the JWT token as a secure, HttpOnly cookie
	c.SetCookie(
		"jwt",                                 // Cookie name
		token,                                 // Cookie value (the JWT token)
		int(expire.Sub(time.Now()).Seconds()), // Cookie max age in seconds
		"/",                                   // Path
		"",                                    // Domain (leave empty for default/current domain)
		true,                                  // Secure (true ensures it is sent over HTTPS only)
		true,                                  // HttpOnly (true prevents JavaScript from accessing the cookie)
	)

	// Optionally, you can return a response without the token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"expire":  expire.Format(time.RFC3339), // Optionally return the expiration time
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the JWT cookie by setting its expiration to a past date
	c.SetCookie(
		"jwt", // Cookie name
		"",    // Empty value
		-1,    // Expiration (set to past to delete the cookie)
		"/",   // Path
		"",    // Domain
		true,  // Secure
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
