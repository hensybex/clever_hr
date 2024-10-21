// model/resume_analysis_result.go

package model

import (
	"time"
)

type AnalysisStatus string

const (
	Pending          AnalysisStatus = "pending"
	InProgress       AnalysisStatus = "in_progress"
	AnalysisFinished AnalysisStatus = "finished"
)

type AnalysisField struct {
	Overview string `json:"overview" gorm:"type:text"`
	Score    int    `json:"score"`
}

type ResumeAnalysisResult struct {
	ID             uint           `gorm:"primaryKey"`
	ResumeID       uint           `gorm:"index"`
	Resume         Resume         `gorm:"foreignKey:ResumeID"`
	AnalysisStatus AnalysisStatus `gorm:"type:varchar(20)"`

	// Individual fields from analysis
	ProfessionalSummaryAndCareerNarrative          AnalysisField `gorm:"embedded;embeddedPrefix:psn_"`
	WorkExperienceAndImpact                        AnalysisField `gorm:"embedded;embeddedPrefix:wei_"`
	EducationAndContinuousLearning                 AnalysisField `gorm:"embedded;embeddedPrefix:ecl_"`
	SkillsAndTechnologicalProficiency              AnalysisField `gorm:"embedded;embeddedPrefix:stp_"`
	SoftSkillsAndEmotionalIntelligence             AnalysisField `gorm:"embedded;embeddedPrefix:sse_"`
	LeadershipInnovationAndProblemSolvingPotential AnalysisField `gorm:"embedded;embeddedPrefix:lip_"`
	CulturalFitAndValueAlignment                   AnalysisField `gorm:"embedded;embeddedPrefix:cfa_"`
	AdaptabilityResilienceAndWorkEthic             AnalysisField `gorm:"embedded;embeddedPrefix:arw_"`
	LanguageProficiencyAndCommunication            AnalysisField `gorm:"embedded;embeddedPrefix:lpc_"`
	ProfessionalAffiliationsAndCommunity           AnalysisField `gorm:"embedded;embeddedPrefix:pac_"`

	// Foreign key to track GPT call details
	GPTCallID *uint    `gorm:"index"`
	GPTCall   *GPTCall `gorm:"foreignKey:GPTCallID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
