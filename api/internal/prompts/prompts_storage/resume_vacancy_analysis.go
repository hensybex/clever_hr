// internal/prompts/prompts_storage/resume_vacancy_analysis.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// ResumeVacancyAnalysisData contains the data for comprehensive resume and vacancy analysis
type ResumeVacancyAnalysisData struct {
	ResumeText   string
	VacancyText  string
	VacancyTitle string
}

// ToPassedData transforms ResumeVacancyAnalysisData into PassedData
func (r ResumeVacancyAnalysisData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "Formatted resume text", Content: r.ResumeText},
		{Name: "Vacancy Text", Description: "Description of the vacancy", Content: r.VacancyText},
		{Name: "Vacancy Title", Description: "Title of the vacancy", Content: r.VacancyTitle},
	}
}

// ResumeVacancyAnalysisPrompt is the prompt configuration for analyzing the match between a resume and a vacancy
var ResumeVacancyAnalysisPrompt = types.Prompt{
	BasePrompt: "You are a highly advanced assistant tasked with providing an in-depth analysis of how well a candidate's resume aligns with a specific job vacancy.",
	BaseTaskDesc: `**Your Task:**

1. **Analyze the provided "Resume Text" and "Vacancy Text".**

2. **For each of the following 10 categories, you must:**
   - **Provide a concise, insightful text overview** summarizing how the candidate's qualifications, experience, and skills align with the requirements of the vacancy. The overviews should be written in Russian language.
   - **Assign a score between 1 and 10** indicating the degree of alignment or suitability in that category, where 1 represents the lowest level and 10 the highest.

3. **Categories to analyze:**

   **1. Relevant Work Experience**
   - Assess how the candidate's past job roles, responsibilities, and achievements match the specific requirements and responsibilities of the vacancy.

   **2. Technical Skills and Proficiencies**
   - Evaluate the candidate's technical skills, tools, and technologies listed in their resume against those required in the vacancy.

   **3. Education and Certifications**
   - Compare the candidate's educational background and certifications with the educational requirements of the vacancy.

   **4. Soft Skills and Cultural Fit**
   - Determine how the candidate's soft skills (e.g., communication, teamwork, leadership) align with the company culture and the soft skills emphasized in the vacancy.

   **5. Language and Communication Skills**
   - Evaluate the candidate's language proficiency and communication abilities in relation to the requirements of the vacancy.

   **6. Problem-Solving and Analytical Abilities**
   - Analyze any evidence of the candidate's problem-solving skills and analytical abilities that are relevant to the vacancy.

   **7. Adaptability and Learning Capacity**
   - Assess the candidate's ability to adapt to new environments and learn new skills as required by the vacancy.

   **8. Leadership and Management Experience**
   - Evaluate any leadership or management experience the candidate has in relation to the leadership expectations of the vacancy.

   **9. Motivation and Career Objectives**
   - Determine how the candidate's career goals and motivations align with the opportunities provided by the vacancy.

   **10. Additional Qualifications and Value-Adds**
   - Identify any additional qualifications, experiences, or attributes the candidate possesses that would add value to the role beyond the basic requirements.

4. **Return a JSON object** with the following structure:
   {
       "RelevantWorkExperience": {
           "Overview": "string",
           "Score": "integer"
       },
       "TechnicalSkillsAndProficiencies": {
           "Overview": "string",
           "Score": "integer"
       },
       "EducationAndCertifications": {
           "Overview": "string",
           "Score": "integer"
       },
       "SoftSkillsAndCulturalFit": {
           "Overview": "string",
           "Score": "integer"
       },
       "LanguageAndCommunicationSkills": {
           "Overview": "string",
           "Score": "integer"
       },
       "ProblemSolvingAndAnalyticalAbilities": {
           "Overview": "string",
           "Score": "integer"
       },
       "AdaptabilityAndLearningCapacity": {
           "Overview": "string",
           "Score": "integer"
       },
       "LeadershipAndManagementExperience": {
           "Overview": "string",
           "Score": "integer"
       },
       "MotivationAndCareerObjectives": {
           "Overview": "string",
           "Score": "integer"
       },
       "AdditionalQualificationsAndValueAdds": {
           "Overview": "string",
           "Score": "integer"
       }
   }
The "Score" for each category should be a number between 1 and 10, where 1 indicates the lowest alignment or suitability, and 10 indicates the highest. Remember that Overviews should be written in Russian language.

Respond with the JSON object only.`,
	UndefinedJSONOutputs: false,
}
