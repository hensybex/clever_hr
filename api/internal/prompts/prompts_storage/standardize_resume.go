// internal/prompts/prompts_storage/standardize_resume.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// ResumeStandardizationData contains the data for resume standardization
type ResumeStandardizationData struct {
	ResumeText string
}

// ToPassedData transforms ResumeStandardizationData into PassedData
func (c ResumeStandardizationData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "The full text extracted from the candidate's resume", Content: c.ResumeText},
	}
}

var ResumeStandardizationPrompt = types.Prompt{
	BasePrompt: "You are an expert assistant specializing in extracting key information from resumes.",
	BaseTaskDesc: `Your task is to analyze the provided "Resume Text" and extract the most critical information into a concise, standardized format.

**Instructions**:

- **Extract** the following key aspects:

  1. **Personal Information**:
     - Full Name
     - Email
     - Phone

  2. **Professional Summary**:
     - A brief overview of the candidate's experience and key skills.

  3. **Key Skills**:
     - Top 5 Technical Skills
     - Top 3 Soft Skills

  4. **Work Experience**:
     - Most recent Job Title
     - Company Name
     - Duration (Start and End Dates)

  5. **Education**:
     - Highest Degree
     - Institution

- **Output Format**:

  - The output must be well-organized and prepared to be passed to another prompt for further processing or review, with the extracted information structured as per the sections above.
  - If any section is not present in the resume, include the section with an empty value.

- **Language**:

  - All extracted text should be in the same language as the "Resume Text" provided.
  - Do not translate or modify the content beyond organizing it into the standardized format.

**Respond with the rewritten resume in plain text, without JSON or other formats.**`,
	UndefinedJSONOutputs: false,
}
