package usecases

import (
// "clever_hr_embeddings/internal/models"
// "clever_hr_embeddings/internal/repository"
// "clever_hr_embeddings/internal/services"
)

type ResumeAnalysisUsecase interface {
	//ExtractSkillsFromResume(resume *models.Resume) (map[string]interface{}, error)
}

type resumeAnalysisUsecaseImpl struct{}

func NewResumeAnalysisUsecase() ResumeAnalysisUsecase {
	return &resumeAnalysisUsecaseImpl{}
}

// ExtractSkillsFromResume calls the LLM function to analyze the resume content and extract skills.
/* func (uc *resumeAnalysisUsecaseImpl) ExtractSkillsFromResume(resume *models.Resume) (map[string]interface{}, error) {
	// Call your LLM function here
	result, err := mistral.NewMistralClient.AnalyzeResumeContent(resume.Content)
	if err != nil {
		return nil, err
	}

	return result, nil
}
*/
