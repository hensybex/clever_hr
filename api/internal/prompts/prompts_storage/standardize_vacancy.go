// prompts/prompts_storage/standardize_vacancy.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// VacancyStandardizationData contains the data for vacancy standardization
type VacancyStandardizationData struct {
	VacancyText string
}

// ToPassedData transforms VacancyStandardizationData into PassedData
func (c VacancyStandardizationData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Vacancy Text", Description: "The full text of the job vacancy", Content: c.VacancyText},
	}
}

var VacancyStandardizationPrompt = types.Prompt{
	BasePrompt: "You are an expert assistant specializing in extracting key information from job vacancies.",
	BaseTaskDesc: `Your task is to analyze the provided "Vacancy Text" and extract the most critical information into a concise, standardized format.

**Instructions**:

- **Extract** the following key aspects:

  1. **Job Title**

  2. **Company Overview**:
     - Brief description of the company.

  3. **Job Summary**:
     - A concise overview of the role and its importance.

  4. **Key Responsibilities**:
     - Top 3 main duties.

  5. **Requirements**:
     - Top 5 required skills or qualifications.

  6. **Preferred Qualifications**:
     - Additional desired skills or qualifications.

  7. **Benefits**:
     - Key benefits offered.


- **Output Format**:

  - The output must be well-organized and prepared to be passed to another prompt for further processing or review, with the extracted information structured as per the sections above.
  - If any section is not present in the resume, include the section with an empty value.

- **Language**:

  - All extracted text should be in the same language as the "Vacancy Text" provided.
  - Do not translate or modify the content beyond organizing it into the standardized format.

**Respond with the rewritten resume in plain text, without JSON or other formats.**`,
	UndefinedJSONOutputs: false,
}
