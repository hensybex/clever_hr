// internal/model/vacancy_resume_match.go

package model

import (
	"time"

	"gorm.io/gorm"
)

type VacancyResumeMatch struct {
	ID        uint    `gorm:"primaryKey"`
	VacancyID uint    `gorm:"index;not null"`
	ResumeID  uint    `gorm:"index;not null"`
	Score     float64 `gorm:"not null"`
	// Analysis fields
	RelevantWorkExperience               AnalysisField `gorm:"embedded;embeddedPrefix:rwe_"`
	TechnicalSkillsAndProficiencies      AnalysisField `gorm:"embedded;embeddedPrefix:tsap_"`
	EducationAndCertifications           AnalysisField `gorm:"embedded;embeddedPrefix:eac_"`
	SoftSkillsAndCulturalFit             AnalysisField `gorm:"embedded;embeddedPrefix:sscf_"`
	LanguageAndCommunicationSkills       AnalysisField `gorm:"embedded;embeddedPrefix:lacs_"`
	ProblemSolvingAndAnalyticalAbilities AnalysisField `gorm:"embedded;embeddedPrefix:psaa_"`
	AdaptabilityAndLearningCapacity      AnalysisField `gorm:"embedded;embeddedPrefix:aalc_"`
	LeadershipAndManagementExperience    AnalysisField `gorm:"embedded;embeddedPrefix:lame_"`
	MotivationAndCareerObjectives        AnalysisField `gorm:"embedded;embeddedPrefix:mco_"`
	AdditionalQualificationsAndValueAdds AnalysisField `gorm:"embedded;embeddedPrefix:aqava_"`
	AnalysisStatus                       string        `gorm:"type:varchar(20);default:'pending'"`

	// Foreign key to track GPT call details (optional)
	GPTCallID *uint    `gorm:"index"`
	GPTCall   *GPTCall `gorm:"foreignKey:GPTCallID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
