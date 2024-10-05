// internal/dtos/create_interview_dto.go

package dtos

// CreateInterviewDTO is the Data Transfer Object for creating an interview
type CreateInterviewDTO struct {
	TgID            uint `json:"tg_id" binding:"required"`
	InterviewTypeID uint `json:"interview_type_id" binding:"required"`
}
