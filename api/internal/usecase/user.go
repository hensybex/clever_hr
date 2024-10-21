// usecase/user.go

package usecase

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/repository"
)

type UserUsecase interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	SwitchUserType(userID uint) error
	//GetCandidates(userID uint) error
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

func (u *userUsecase) SwitchUserType(userID uint) error {
	return u.userRepo.SwitchUserType(userID)
}

/* func (u *userUsecase) GetCandidates(userID uint) error {
	candidates, err := u.userRepo.FindCandidatesByUserID(userID)
	if err != nil {
		return err
	}

	// Perform additional logic here to get the list of candidates uploaded by this user
	// For example, you can fetch the candidates from a separate repository or service

	return nil
}
*/
