// internal/dtos/create_user_dto.go

package dtos

import (
	"clever_hr_api/internal/model"
	"encoding/json"
)

type CreateUserDTO struct {
	TgID     json.Number `json:"tg_id" binding:"required"`
	UserType string      `json:"user_type" binding:"required,oneof=employee candidate"`
}

// ToUserModel converts CreateUserDTO to a User model
func (dto CreateUserDTO) ToUserModel() model.User {
	return model.User{
		TgID:     dto.TgID.String(), // Convert json.Number to string
		UserType: model.UserType(dto.UserType),
	}
}
