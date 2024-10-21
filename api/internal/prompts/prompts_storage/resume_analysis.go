// prompts/prompts_storage/resume_analysis.go

package prompts_storage

import "clever_hr_api/internal/prompts/types"

// ResumeAnalysisData contains the data for comprehensive resume analysis
type ResumeAnalysisData struct {
	ResumeText string
}

// ToPassedData transforms ResumeRewritionData into PassedData
func (r ResumeAnalysisData) ToPassedData() []types.PassedData {
	return []types.PassedData{
		{Name: "Resume Text", Description: "Formatted resume text", Content: r.ResumeText},
	}
}

// ResumeAnalysisPrompt is the prompt configuration for resume analysis
var ResumeAnalysisPrompt = types.Prompt{
	BasePrompt: "You are a highly advanced assistant tasked with providing in-depth resume analysis to help HR professionals make well-informed hiring decisions.",
	BaseTaskDesc: `**Your Task:**

1. **Analyze the provided "Resume Text".**

2. **For each of the following 10 categories, you must:**
   - **Provide a concise, insightful text overview** summarizing the candidate's strengths and relevant details specific to that category. Text overviews should be written in Russian language.
   - **Assign a score between 1 and 10** indicating the candidate's proficiency, suitability, or alignment in that category, where 1 represents the lowest level and 10 the highest.

3. **Categories to analyze:**

   **1. Professional Summary and Career Narrative**
   - Evaluate how well the candidate presents their professional journey. Assess the clarity of their career objectives, the logical progression of roles, and their ability to communicate this effectively. Are they able to weave a clear, strategic narrative about their career path?

   **2. Work Experience and Impact**
   - Analyze the candidate’s job titles, responsibilities, and achievements across past roles. Focus on measurable contributions, such as revenue growth, cost savings, or project completions. Look for evidence of increasing responsibility and a pattern of consistent performance.

   **3. Education and Continuous Learning**
   - Examine the candidate’s formal education and any additional certifications, courses, or self-directed learning. Assess their commitment to professional development and how their education aligns with the role they are seeking.

   **4. Skills and Technological Proficiency**
   - Identify the technical and non-technical skills the candidate possesses. Determine the depth of their proficiency with the tools, technologies, or methodologies relevant to the job. Evaluate if they demonstrate an ability to quickly learn and adapt to new technologies.

   **5. Soft Skills and Emotional Intelligence**
   - Evaluate the candidate’s interpersonal skills, such as teamwork, leadership, communication, and adaptability. Look for evidence of emotional intelligence, empathy, and the ability to navigate complex social dynamics in a professional setting.

   **6. Leadership, Innovation, and Problem-Solving Potential**
   - Assess the candidate’s demonstrated leadership abilities, innovative thinking, and their capacity for creative problem-solving. Look for examples where they led teams, spearheaded initiatives, or solved complex problems that led to business success.

   **7. Cultural Fit and Value Alignment**
   - Analyze how well the candidate’s values, ethics, and cultural perspectives align with the organization’s mission. Consider their past experiences in terms of teamwork, diversity, and ethical decision-making, and determine if their professional ethos resonates with the company’s core values.

   **8. Adaptability, Resilience, and Work Ethic**
   - Assess how well the candidate adapts to change, overcomes challenges, and perseveres in the face of adversity. Look for instances that highlight their resilience and strong work ethic, especially in difficult or fast-paced environments.

   **9. Language Proficiency and Communication Skills**
   - Examine the candidate’s communication abilities, including written and verbal communication skills. If applicable, evaluate their proficiency in multiple languages and how well they articulate ideas and collaborate with others.

   **10. Professional Affiliations and Community Engagement**
   - Assess the candidate’s involvement in professional organizations, volunteer work, or community initiatives. Evaluate how these activities demonstrate their commitment to professional growth, networking, and social responsibility.

4. **Return a JSON object** with the following structure:
   {
       "ProfessionalSummaryAndCareerNarrative": {
           "Overview": "string",
           "Score": "integer"
       },
       "WorkExperienceAndImpact": {
           "Overview": "string",
           "Score": "integer"
       },
       "EducationAndContinuousLearning": {
           "Overview": "string",
           "Score": "integer"
       },
       "SkillsAndTechnologicalProficiency": {
           "Overview": "string",
           "Score": "integer"
       },
       "SoftSkillsAndEmotionalIntelligence": {
           "Overview": "string",
           "Score": "integer"
       },
       "LeadershipInnovationAndProblemSolvingPotential": {
           "Overview": "string",
           "Score": "integer"
       },
       "CulturalFitAndValueAlignment": {
           "Overview": "string",
           "Score": "integer"
       },
       "AdaptabilityResilienceAndWorkEthic": {
           "Overview": "string",
           "Score": "integer"
       },
       "LanguageProficiencyAndCommunicationSkills": {
           "Overview": "string",
           "Score": "integer"
       },
       "ProfessionalAffiliationsAndCommunityEngagement": {
           "Overview": "string",
           "Score": "integer"
       }
   }
The "Score" for each category should be a number between 1 and 10, where 1 indicates the lowest proficiency or suitability, and 10 indicates the highest. Remember that Overviews should be written in Russian language.
Respond with the JSON object only.
`,
	UndefinedJSONOutputs: false,
}
