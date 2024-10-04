// internal/dtos/create_interview_dto.go

package dtos

import "clever_hr_api/internal/model"

// CreateInterviewDTO is the Data Transfer Object for creating an interview
type CreateInterviewDTO struct {
	ResumeID        uint `json:"resume_id" binding:"required"`
	InterviewTypeID uint `json:"interview_type_id" binding:"required"`
}

// ToInterviewModel converts CreateInterviewDTO to the Interview model
func (dto CreateInterviewDTO) ToInterviewModel() model.Interview {
	return model.Interview{
		ResumeID:        dto.ResumeID,
		InterviewTypeID: dto.InterviewTypeID,
		Status:          model.InterviewStatus(model.InProgress),
	}
}
