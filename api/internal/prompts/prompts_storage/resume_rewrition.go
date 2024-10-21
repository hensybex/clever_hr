// prompts/prompts_storage/resume_rewrition.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// ResumeRewritionData contains the data for resume analysis
type ResumeRewritionData struct {
	ResumeText string
}

// ToPassedData transforms ResumeRewritionData into PassedData
func (r ResumeRewritionData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "The full text extracted from the candidate's resume", Content: r.ResumeText},
	}
}

var ResumeRewritionPrompt = types.Prompt{
	BasePrompt: "You are an assistant tasked with reviewing resumes.",
	BaseTaskDesc: `**Your Task:**

1. Take the "Resume Text" provided.

2. Rewrite the resume text in a more structured, human-readable format, without omitting any information. Ensure that all sections, such as education, work experience, skills, and other relevant details, are clearly delineated.

3. The output must be well-organized and prepared to be passed to another prompt for further processing or review.

**Respond with the rewritten resume in plain text, without JSON or other formats.**`,
	UndefinedJSONOutputs: false,
}
