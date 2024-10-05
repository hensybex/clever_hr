// internal/usecase/user_usecase.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
)

type UserUsecase interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByTgID(tgID string) (*model.User, error)
	SwitchUserType(userID uint) error
	GetCandidatesByTGID(userID string) ([]model.Candidate, error)
	GetUserRoleByTgID(tgID string) (string, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) CreateUser(user *model.User) error {
	return u.userRepo.CreateUser(user)
}

func (u *userUsecase) GetUserByID(id uint) (*model.User, error) {
	return u.userRepo.GetUserByID(id)
}

func (u *userUsecase) GetUserByTgID(tgID string) (*model.User, error) {
	return u.userRepo.GetUserByTgID(tgID)
}

func (u *userUsecase) SwitchUserType(userID uint) error {
	return u.userRepo.SwitchUserType(userID)
}

func (u *userUsecase) GetCandidatesByTGID(telegramID string) ([]model.Candidate, error) {
	// Call the repository to get the list of candidates
	// Retrieve the user by their Telegram ID
	user, err := u.userRepo.GetUserByTgID(telegramID)
	if err != nil {
		return nil, nil
	}
	candidates, err := u.userRepo.FindCandidatesByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func (u *userUsecase) GetUserRoleByTgID(tgID string) (string, error) {
	// Call the repository to get the user information
	user, err := u.userRepo.GetUserByTgID(tgID)
	if err != nil {
		return "", err
	}

	// Check if the user exists
	if user == nil {
		return "", nil // This will return a 404 in the handler
	}

	// Return the user's role (employee or candidate)
	return string(user.UserType), nil
}
