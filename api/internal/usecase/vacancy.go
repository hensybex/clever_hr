// usecase/vacancy.go

package usecase

import (
	"clever_hr_api/internal/dtos"
	"clever_hr_api/internal/mistral"
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/prompts"
	"clever_hr_api/internal/prompts/prompts_storage"
	"clever_hr_api/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type VacancyUsecase interface {
	//CreateVacancy(vacancy *model.Vacancy) error
	UploadVacancy(vacancy *model.Vacancy, conn *websocket.Conn) error
	GetVacancyByID(id uint) (*model.Vacancy, error)
	GetAllVacancies() ([]model.Vacancy, error)
	UpdateVacancyStatus(id uint, status string) error
	UpdateAllVacancyEmbeddings() error
}

type vacancyUsecase struct {
	vacancyRepo        repository.VacancyRepository
	embeddingRepo      repository.EmbeddingRepository
	jobGroupRepo       repository.JobGroupRepository
	specializationRepo repository.SpecializationRepository
	qualificationRepo  repository.QualificationRepository
	mistralService     mistral.MistralService
	matchUsecase       VacancyResumeMatchUsecase
}

func NewVacancyUsecase(
	vacancyRepo repository.VacancyRepository,
	embeddingRepo repository.EmbeddingRepository,
	jobGroupRepo repository.JobGroupRepository,
	specializationRepo repository.SpecializationRepository,
	qualificationRepo repository.QualificationRepository,
	mistralService mistral.MistralService,
	matchUsecase VacancyResumeMatchUsecase,
) VacancyUsecase {
	return &vacancyUsecase{
		vacancyRepo,
		embeddingRepo,
		jobGroupRepo,
		specializationRepo,
		qualificationRepo,
		mistralService,
		matchUsecase,
	}
}

type GroupDetectionVacancyResponse struct {
	Analysis  string `json:"Analysis"`
	JobTypeID uint   `json:"JobTypeID"`
	Title     string `json:"Title"`
}

/* func (u *vacancyUsecase) UploadVacancy(vacancy *model.Vacancy) error {
	// Log the start of the function and input parameters
	log.Printf("Starting UploadVacancy")

	// Get Groups Data
	jobGroups, err := u.jobGroupRepo.GetAll()
	if err != nil {
		log.Printf("Error getting all job groups. error: %v", err)
		return err
	}
	var jobGroupsString string

	for _, group := range jobGroups {
		jobGroupsString += fmt.Sprintf("Group name: %s, id: %d; ", group.TitleEn, group.ID)
	}

	jobGroupsString = strings.TrimSuffix(jobGroupsString, "; ")

	// Step 1: Construct LLM prompt for Group Detection
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	groupDetectionData := prompts_storage.GroupDetectionVacancyData{VacancyDescription: vacancy.Description, JobTypes: jobGroupsString}
	prompt, err := pc.GetPrompt(allPrompts.GroupDetectionVacancyPrompt, groupDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("GroupDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the vacancy
	llmResponse, _, err := u.mistralService.CallMistral(prompt, true, mistral.Nemo, "GroupDetectionVacancy", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var groupDetectionResponse GroupDetectionVacancyResponse
	if err := json.Unmarshal([]byte(llmResponse), &groupDetectionResponse); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Step 2: Construct LLM prompt for Specialization Detection
	// Obtain Job Group
	jobGroup, err := u.jobGroupRepo.GetByID(int(groupDetectionResponse.JobTypeID))
	if err != nil {
		log.Printf("Error getting job group by ID. error: %v", err)
		return err
	}

	// Obtain correct specialization
	specializations, err := u.specializationRepo.GetAllByGroupID(int(groupDetectionResponse.JobTypeID))
	if err != nil {
		log.Printf("Error getting all job groups. error: %v", err)
		return err
	}
	var specializationsString string

	for _, specialization := range specializations {
		specializationsString += fmt.Sprintf("Specialization name: %s, id: %d; ", specialization.TitleEn, specialization.ID)
	}

	specializationsString = strings.TrimSuffix(specializationsString, "; ")

	specializationDetectionData := prompts_storage.SpecializationDetectionVacancyData{VacancyDescription: vacancy.Description, JobGroupType: jobGroup.TitleEn, SpecializationTypes: specializationsString}
	specializationPrompt, err := pc.GetPrompt(allPrompts.SpecializationDetectionVacancyPrompt, specializationDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("SpecializationDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the vacancy
	llmResponse, _, err = u.mistralService.CallMistral(specializationPrompt, true, mistral.Nemo, "SpecializationDetectionVacancy", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var specializationResponse SpecializationDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &specializationResponse); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Step 3: Construct LLM prompt for Qualification Detection
	// Obtain Specialization
	specialization, err := u.specializationRepo.GetByID(int(specializationResponse.SpecializationTypeID))
	if err != nil {
		log.Printf("Error getting specialization by ID. error: %v", err)
		return err
	}

	// Obtain qualifications
	qualifications, err := u.qualificationRepo.GetAll()
	if err != nil {
		log.Printf("Error getting all job groups. error: %v", err)
		return err
	}
	var qualificationsString string

	for _, qualification := range qualifications {
		qualificationsString += fmt.Sprintf("Qualification name: %s, id: %d; ", qualification.TitleEn, qualification.ID)
	}

	qualificationsString = strings.TrimSuffix(qualificationsString, "; ")

	qualificationDetectionData := prompts_storage.QualificationDetectionVacancyData{VacancyDescription: vacancy.Description, JobGroupType: jobGroup.TitleEn, SpecializationType: specialization.TitleEn, QualificationTypes: qualificationsString}
	qualificationPrompt, err := pc.GetPrompt(allPrompts.QualificationDetectionVacancyPrompt, qualificationDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("QualificationDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the vacancy
	llmResponse, _, err = u.mistralService.CallMistral(qualificationPrompt, true, mistral.Nemo, "QualificationDetectionVacancy", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var qualificationResponse QualificationDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &qualificationResponse); err != nil {
		log.Printf("Error parsing LLM response. error: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Save the extracted text and rewritten vacancy in the database
	vacancy.JobGroupID = &groupDetectionResponse.JobTypeID
	vacancy.SpecializationID = &specializationResponse.SpecializationTypeID
	vacancy.QualificationID = &qualificationResponse.QualificationTypeID
	vacancy.Title = groupDetectionResponse.Title
	log.Printf("Saving vacancy to the database. Resume ID: %v", vacancy.ID)
	if err := u.vacancyRepo.CreateOne(vacancy); err != nil {
		log.Printf("Error saving vacancy. error: %v", err)
		return fmt.Errorf("error saving vacancy: %v", err)
	}

	log.Printf("Starting generating embedding")
	embedding, err := u.mistralService.GenerateEmbedding(vacancy.Description)
	if err != nil {
		return err
	}

	// Convert the []float64 embedding to []float32 for pgvector
	embeddingFloat32 := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingFloat32[i] = float32(v)
	}

	// Save the embedding as pgvector.Vector
	vacancyEmbedding := &model.VacancyEmbedding{
		VacancyID: vacancy.ID,
		Embedding: pgvector.NewVector(embeddingFloat32), // Correct conversion
	}

	if err := u.embeddingRepo.CreateVacancyEmbedding(vacancyEmbedding); err != nil {
		return err
	}

	// Log the successful completion of the function
	log.Printf("UploadResume completed successfully for Resume ID: %v", vacancy.ID)

	return nil
} */

func (u *vacancyUsecase) UploadVacancy(vacancy *model.Vacancy, conn *websocket.Conn) error {
	log.Printf("Starting UploadVacancy for Vacancy ID: %v", vacancy.ID)

	// Send status update
	conn.WriteJSON(dtos.ServerMessage{Status: "Fetching job groups"})

	// Get Groups Data
	jobGroups, err := u.jobGroupRepo.GetAll()
	if err != nil {
		log.Printf("Error getting all job groups: %v", err)
		return err
	}
	var jobGroupsString string

	for _, group := range jobGroups {
		jobGroupsString += fmt.Sprintf("Group name: %s, id: %d; ", group.TitleEn, group.ID)
	}
	jobGroupsString = strings.TrimSuffix(jobGroupsString, "; ")

	// Step 0: Construct LLM prompt for Resume Standarization
	conn.WriteJSON(dtos.ServerMessage{Status: "Standarizing Vacancy Description"})
	allPrompts := prompts.NewPrompts()
	pc := prompts.NewPromptConstructor()
	standarizationData := prompts_storage.VacancyStandardizationData{VacancyText: vacancy.Description}
	prompt, err := pc.GetPrompt(allPrompts.VacancyStandardizationPrompt, standarizationData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt. error: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("GroupDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the resume
	llmResponse, _, err := u.mistralService.CallMistral(prompt, false, mistral.Largest, "VacancyStandarization", 0)
	if err != nil {
		log.Printf("Error calling LLM service. error: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	standarizedVacancyText := llmResponse

	// Step 1: Construct LLM prompt for Group Detection
	conn.WriteJSON(dtos.ServerMessage{Status: "Constructing group detection prompt"})
	groupDetectionData := prompts_storage.GroupDetectionVacancyData{VacancyDescription: standarizedVacancyText, JobTypes: jobGroupsString}
	prompt, err = pc.GetPrompt(allPrompts.GroupDetectionVacancyPrompt, groupDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating LLM prompt: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}
	log.Printf("GroupDetectionData prompt constructed successfully.")

	// Call LLM to extract candidate details and rewrite the vacancy
	conn.WriteJSON(dtos.ServerMessage{Status: "Calling LLM for group detection"})
	llmResponse, _, err = u.mistralService.CallMistral(prompt, true, mistral.Nemo, "GroupDetectionVacancy", 0)
	if err != nil {
		log.Printf("Error calling LLM service for group detection: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	// Parse LLM response into structured data
	var groupDetectionResponse GroupDetectionVacancyResponse
	if err := json.Unmarshal([]byte(llmResponse), &groupDetectionResponse); err != nil {
		log.Printf("Error parsing LLM response for group detection: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Step 2: Construct LLM prompt for Specialization Detection
	conn.WriteJSON(dtos.ServerMessage{Status: "Fetching specialization data"})
	jobGroup, err := u.jobGroupRepo.GetByID(int(groupDetectionResponse.JobTypeID))
	if err != nil {
		log.Printf("Error getting job group by ID: %v", err)
		return err
	}

	specializations, err := u.specializationRepo.GetAllByGroupID(int(groupDetectionResponse.JobTypeID))
	if err != nil {
		log.Printf("Error getting all specializations: %v", err)
		return err
	}

	var specializationsString string
	for _, specialization := range specializations {
		specializationsString += fmt.Sprintf("Specialization name: %s, id: %d; ", specialization.TitleEn, specialization.ID)
	}
	specializationsString = strings.TrimSuffix(specializationsString, "; ")

	specializationDetectionData := prompts_storage.SpecializationDetectionVacancyData{VacancyDescription: standarizedVacancyText, JobGroupType: jobGroup.TitleEn, SpecializationTypes: specializationsString}
	specializationPrompt, err := pc.GetPrompt(allPrompts.SpecializationDetectionVacancyPrompt, specializationDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating specialization prompt: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}

	conn.WriteJSON(dtos.ServerMessage{Status: "Calling LLM for specialization detection"})
	llmResponse, _, err = u.mistralService.CallMistral(specializationPrompt, true, mistral.Nemo, "SpecializationDetectionVacancy", 0)
	if err != nil {
		log.Printf("Error calling LLM service for specialization detection: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	var specializationResponse SpecializationDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &specializationResponse); err != nil {
		log.Printf("Error parsing LLM response for specialization detection: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Step 3: Construct LLM prompt for Qualification Detection
	conn.WriteJSON(dtos.ServerMessage{Status: "Fetching qualification data"})
	specialization, err := u.specializationRepo.GetByID(int(specializationResponse.SpecializationTypeID))
	if err != nil {
		log.Printf("Error getting specialization by ID: %v", err)
		return err
	}

	qualifications, err := u.qualificationRepo.GetAll()
	if err != nil {
		log.Printf("Error getting all qualifications: %v", err)
		return err
	}

	var qualificationsString string
	for _, qualification := range qualifications {
		qualificationsString += fmt.Sprintf("Qualification name: %s, id: %d; ", qualification.TitleEn, qualification.ID)
	}
	qualificationsString = strings.TrimSuffix(qualificationsString, "; ")

	qualificationDetectionData := prompts_storage.QualificationDetectionVacancyData{VacancyDescription: standarizedVacancyText, JobGroupType: jobGroup.TitleEn, SpecializationType: specialization.TitleEn, QualificationTypes: qualificationsString}
	qualificationPrompt, err := pc.GetPrompt(allPrompts.QualificationDetectionVacancyPrompt, qualificationDetectionData, "", false)
	if err != nil {
		log.Printf("Error generating qualification prompt: %v", err)
		return fmt.Errorf("error generating LLM prompt: %v", err)
	}

	conn.WriteJSON(dtos.ServerMessage{Status: "Calling LLM for qualification detection"})
	llmResponse, _, err = u.mistralService.CallMistral(qualificationPrompt, true, mistral.Nemo, "QualificationDetectionVacancy", 0)
	if err != nil {
		log.Printf("Error calling LLM service for qualification detection: %v", err)
		return fmt.Errorf("error calling LLM: %v", err)
	}

	var qualificationResponse QualificationDetectionResponse
	if err := json.Unmarshal([]byte(llmResponse), &qualificationResponse); err != nil {
		log.Printf("Error parsing LLM response for qualification detection: %v", err)
		return fmt.Errorf("error parsing LLM response: %v", err)
	}

	// Save the extracted data in the vacancy
	conn.WriteJSON(dtos.ServerMessage{Status: "Saving vacancy data"})
	vacancy.JobGroupID = &groupDetectionResponse.JobTypeID
	vacancy.SpecializationID = &specializationResponse.SpecializationTypeID
	vacancy.QualificationID = &qualificationResponse.QualificationTypeID
	vacancy.Title = groupDetectionResponse.Title
	vacancy.StandarizedText = standarizedVacancyText

	log.Printf("Saving vacancy to the database. Vacancy ID: %v", vacancy.ID)
	if err := u.vacancyRepo.CreateOne(vacancy); err != nil {
		log.Printf("Error saving vacancy. error: %v", err)
		return fmt.Errorf("error saving vacancy: %v", err)
	}

	// Embedding generation
	conn.WriteJSON(dtos.ServerMessage{Status: "Generating embedding for vacancy"})
	embedding, err := u.mistralService.GenerateEmbedding(vacancy.StandarizedText)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		return err
	}

	log.Printf("Embedding generated successfully for Vacancy ID: %v", vacancy.ID)

	// Convert the []float64 embedding to []float32
	embeddingFloat32 := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingFloat32[i] = float32(v)
	}

	// Save the embedding
	vacancyEmbedding := &model.VacancyEmbedding{
		VacancyID: vacancy.ID,
		Embedding: embeddingFloat32, // Store the []float32 directly
	}

	log.Printf("Saving embedding for Vacancy ID: %v", vacancy.ID)

	// Pass both vacancyID and embedding to the repository
	if err := u.embeddingRepo.CreateVacancyEmbedding(vacancyEmbedding.VacancyID, vacancyEmbedding.Embedding); err != nil {
		log.Printf("Error saving vacancy embedding: %v", err)
		return err
	}

	// Detailed logs after embedding generation
	conn.WriteJSON(dtos.ServerMessage{Status: "Matching vacancy with resumes"})
	log.Printf("Starting resume matching process for Vacancy ID: %v", vacancy.ID)

	// Call MatchVacancyWithResumes and pass the WebSocket connection
	if err := u.matchUsecase.MatchVacancyWithResumes(vacancy.ID, conn); err != nil {
		log.Printf("Error matching vacancy with resumes: %v", err)
		return err
	}

	// Send final completion status
	conn.WriteJSON(dtos.ServerMessage{Status: "Vacancy processing and matching completed"})
	log.Printf("UploadVacancy and matching completed successfully for Vacancy ID: %v", vacancy.ID)

	return nil
}

func (uc *vacancyUsecase) GetVacancyByID(id uint) (*model.Vacancy, error) {
	return uc.vacancyRepo.GetByID(id)
}

func (uc *vacancyUsecase) GetAllVacancies() ([]model.Vacancy, error) {
	return uc.vacancyRepo.GetAll()
}

func (uc *vacancyUsecase) UpdateVacancyStatus(id uint, status string) error {
	return uc.vacancyRepo.UpdateStatus(id, status)
}

func (u *vacancyUsecase) UpdateAllVacancyEmbeddings() error {
	vacancies, err := u.vacancyRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to retrieve vacancies: %w", err)
	}

	for _, vacancy := range vacancies {
		log.Printf("Processing Vacancy ID: %v, StandarizedText: %v", vacancy.ID, vacancy.StandarizedText)

		// Check if embedding exists
		_, err := u.embeddingRepo.GetVacancyEmbedding(vacancy.ID)
		if err == nil {
			log.Printf("Embedding already exists for Vacancy ID: %v, skipping...", vacancy.ID)
			continue
		}

		// Generate embedding
		embedding, err := u.mistralService.GenerateEmbedding(vacancy.StandarizedText)
		if err != nil {
			log.Printf("Error generating embedding for Vacancy ID: %v, error: %v", vacancy.ID, err)
			continue
		}

		// Log the generated embedding
		log.Printf("Generated embedding for Vacancy ID: %v, embedding: %v", vacancy.ID, embedding)

		// Convert embedding from []float64 to []float32
		embeddingFloat32 := make([]float32, len(embedding))
		for i, v := range embedding {
			embeddingFloat32[i] = float32(v)
		}

		// Save embedding
		if err := u.embeddingRepo.CreateVacancyEmbedding(vacancy.ID, embeddingFloat32); err != nil {
			log.Printf("Error saving embedding for Vacancy ID: %v, error: %v", vacancy.ID, err)
			continue
		}

		log.Printf("Successfully updated embedding for Vacancy ID: %v", vacancy.ID)
	}
	return nil
}
