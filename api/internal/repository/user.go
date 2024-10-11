// internal/repository/user.go

package repository

import (
	"clever_hr_api/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	SwitchUserType(userID uint) error
	FindCandidatesByUserID(userID uint) ([]model.Candidate, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) SwitchUserType(userID uint) error {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Toggle user type between 'employee' and 'candidate'
	if user.UserType == model.Employee {
		user.UserType = model.UserTypeCandidate
	} else {
		user.UserType = model.Employee
	}

	// Update the user type in the database
	return r.db.Save(&user).Error
}

func (r *userRepository) FindCandidatesByUserID(userID uint) ([]model.Candidate, error) {
	var candidates []model.Candidate

	// Query the database to find candidates where UploadedByUserID is the given userID
	if err := r.db.Where("uploaded_by_user_id = ?", userID).Find(&candidates).Error; err != nil {
		return nil, err
	}

	return candidates, nil
}
