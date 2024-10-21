// model/resume.go

package model

import (
	category_model "clever_hr_api/internal/model/categories"
	"time"
)

type Resume struct {
	ID              uint `gorm:"primaryKey"`
	UploadedFrom    string
	PDFPath         string
	RawText         string
	CleanText       string
	StandarizedText string

	JobGroupID       *uint `gorm:"index"`
	SpecializationID *uint `gorm:"index"`
	QualificationID  *uint `gorm:"index"`
	FullName         *string
	GPTCallID        *uint                          `gorm:"index"`
	JobGroup         *category_model.JobGroup       `gorm:"foreignKey:JobGroupID"`
	Specialization   *category_model.Specialization `gorm:"foreignKey:SpecializationID"`
	Qualification    *category_model.Qualification  `gorm:"foreignKey:QualificationID"`
	GPTCall          *GPTCall                       `gorm:"foreignKey:GPTCallID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
