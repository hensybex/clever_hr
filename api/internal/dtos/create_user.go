// dtos/create_user.go

package dtos

import (
	"clever_hr_api/internal/model"
)

type CreateUserDTO struct {
	UserID   uint   `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=employee candidate"`
}

// ToUserModel converts CreateUserDTO to a User model
func (dto CreateUserDTO) ToUserModel() model.User {
	return model.User{
		ID:       dto.UserID, // Convert json.Number to string
		UserType: model.UserType(dto.UserType),
	}
}
