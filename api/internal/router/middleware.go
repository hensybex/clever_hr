// internal/router/middleware.go

package router

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/service"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware - create middleware for JWT
func AuthMiddleware(authService service.AuthService) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Key:         []byte("your_secret_key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					"id":       v.ID,
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			user, err := authService.Authenticate(loginVals.Username, loginVals.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				ID:       uint(claims["id"].(float64)),
				Username: claims["username"].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.User); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
