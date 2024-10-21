// prompts/prompts_storage/resume_initial_analysis.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// ResumeInitialAnalysisData contains the data for resume analysis and candidate extraction
type ResumeInitialAnalysisData struct {
	ResumeText string
}

// ToPassedData transforms ResumeInitialAnalysisData into PassedData
func (c ResumeInitialAnalysisData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "The full text extracted from the candidate's resume", Content: c.ResumeText},
	}
}

var ResumeInitialAnalysisPrompt = types.Prompt{
	BasePrompt: "You are an assistant tasked with analyzing resumes and extracting candidate details.",
	BaseTaskDesc: `**Your Task:**

1. Take the "Resume Text" provided.

2. Rewrite the resume text in a more structured, human-readable format, without omitting any information. Ensure that all sections, such as education, work experience, skills, and other relevant details, are clearly delineated.

3. Extract the following information separately (if available):
   - FullName (if full name is provided, extract it, oterwise omit it)
   - Email (if email is provided, extract it, oterwise omit it)
   - Phone (if phone is provided, extract it, oterwise omit it)
   - BirthDate (if birth date is provided, extract it in format "YYYY-MM-DD", oterwise omit it)
   - TotalYears (if total years of experience is explicitly mentioned, extract it, oterwise omit it)
   - PreferableJob (if mentioned, extract it; if not, suggest based on the resume)


Provide BirthDate and TotalYears only if you are absolutely certain that you found them correctly. BirthDate and TotalYears usually come at the start of the resume. Assume that resumes might be written in Russian, so BirthDate and TotalYears might be near a keyword like "Дата рождения", "Возраст" etc.

4. Return a JSON object with the following fields:
   - "RewrittenResume": The rewritten resume text in plain text format
   - "FullName": Candidate's full name
   - "Email": Candidate's email
   - "Phone": Candidate's phone number (if provided, otherwise "")
   - "BirthDate": Candidate's birth date in format "YYYY-MM-DD" (if available, otherwise "")
   - "TotalYears": Total years of experience as String (if explicitly mentioned)
   - "PreferableJob": The preferable job mentioned or suggested based on resume content

**Respond with the JSON object only.**`,
	UndefinedJSONOutputs: false,
}
