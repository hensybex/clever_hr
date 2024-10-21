// internal/dtos/create_interview.go

package dtos

// CreateInterviewDTO is the Data Transfer Object for creating an interview
type CreateInterviewDTO struct {
	UserID          uint `json:"user_id" binding:"required"`
	InterviewTypeID uint `json:"interview_type_id" binding:"required"`
}
