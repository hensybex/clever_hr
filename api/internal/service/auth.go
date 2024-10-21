// service/auth.go

package service

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(username, password string) (*model.User, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	GenerateToken(user *model.User) (string, time.Time, error)
}

type authService struct {
	userRepository repository.UserRepository
	jwtMiddleware  *jwt.GinJWTMiddleware
}

func NewAuthService(repo repository.UserRepository, jwtMiddleware *jwt.GinJWTMiddleware) AuthService {
	return &authService{userRepository: repo, jwtMiddleware: jwtMiddleware}
}

func (s *authService) Authenticate(username, password string) (*model.User, error) {
	user, err := s.userRepository.FindByUsername(username)
	if err != nil || s.ComparePassword(user.Password, password) != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	return user, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *authService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken generates a JWT token for the authenticated user
func (s *authService) GenerateToken(user *model.User) (string, time.Time, error) {
	token, expire, err := s.jwtMiddleware.TokenGenerator(user)
	return token, expire, err
}
