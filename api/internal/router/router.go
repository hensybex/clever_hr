// router/router.go

package router

import (
	"clever_hr_api/internal/handlers"
	"clever_hr_api/internal/mistral"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/service"
	"clever_hr_api/internal/usecase"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, milvusClient client.Client) *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	resumeRepo := repository.NewResumeRepository(db)
	resumeAnalysisResultRepo := repository.NewResumeAnalysisResultRepository(db)
	gptCallRepo := repository.NewGPTCallRepository(db)
	vacancyRepo := repository.NewVacancyRepository(db)
	jobGroupRepo := repository.NewJobGroupRepository(db)
	specializationRepo := repository.NewSpecializationRepository(db)
	qualificationRepo := repository.NewQualificationRepository(db)
	vacancyResumeMatchRepo := repository.NewVacancyResumeMatchRepository(db)
	embeddingRepo := repository.NewEmbeddingRepository(milvusClient)

	// Initialize services
	mistralService := mistral.NewMistralService(gptCallRepo)
	authService := service.NewAuthService(userRepo, nil)

	// Initialize AuthMiddleware
	authMiddleware, err := AuthMiddleware(authService)
	if err != nil {
		log.Fatal("JWT Error: ", err)
	}

	// Set the jwtMiddleware in AuthService after it's created
	//authService = service.NewAuthService(userRepo, authMiddleware)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	resumeUsecase := usecase.NewResumeUsecase(resumeRepo, embeddingRepo, resumeAnalysisResultRepo, userRepo, jobGroupRepo, specializationRepo, qualificationRepo, *mistralService)
	resumeAnalysisResultUsecase := usecase.NewResumeAnalysisResultUsecase(resumeAnalysisResultRepo, resumeRepo)
	matchUsecase := usecase.NewMatchUsecase(embeddingRepo, vacancyResumeMatchRepo, vacancyRepo, resumeRepo, *mistralService)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo, embeddingRepo, jobGroupRepo, specializationRepo, qualificationRepo, *mistralService, matchUsecase)
	//embeddingUsecase := usecase.NewEmbeddingUsecase(embeddingRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUsecase)
	resumeHandler := handlers.NewResumeHandler(resumeUsecase, userUsecase, resumeAnalysisResultUsecase)
	vacancyHandler := handlers.NewVacancyHandler(vacancyUsecase)
	matchHandler := handlers.NewMatchHandler(matchUsecase)
	authHandler := handlers.NewAuthHandler(authService)
	//embeddingHandler := handlers.NewEmbeddingHandler(embeddingUsecase)

	// Public routes (no auth needed)
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/logout", authHandler.Logout)

	// Define API routes
	api := r.Group("/api")
	// Update the API routes to be more specific
	//api.GET("/embeddings/resume/:resume_id/embedding", embeddingHandler.GetResumeEmbedding)
	//api.GET("/embeddings/vacancy/:vacancy_id/embedding", embeddingHandler.GetVacancyEmbedding)
	//api.GET("/embeddings/vacancy/:vacancy_id/match-resumes", embeddingHandler.FindMatchingResumes)

	api.Use(authMiddleware.MiddlewareFunc())
	{
		// User routes
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:user_id/info", userHandler.GetUser)
		api.PUT("/users/:user_id/switch", userHandler.SwitchUserType)

		// Resume routes
		api.GET("/resumes/:resume_id", resumeHandler.GetResumeByID)
		api.POST("/resumes/upload", resumeHandler.UploadResume)
		api.GET("/resumes/:resume_id/analyze", resumeHandler.RunResumeAnalysis)
		api.GET("/resumes/:resume_id/analysis-result", resumeHandler.GetResumeAnalysisResult)

		// Vacancy routes
		//api.POST("/vacancies/upload", vacancyHandler.UploadVacancy)
		api.GET("/vacancies", vacancyHandler.GetVacancies)
		api.GET("/vacancies/:id", vacancyHandler.GetVacancyByID)
		api.PUT("/vacancies/:id/status", vacancyHandler.UpdateVacancyStatus)

		// WebSocket route for matching vacancy with resumes
		api.GET("/api/vacancies/:vac_id/match", matchHandler.MatchVacancyWithResumes)
		// WebSocket route for uploading a vacancy and matching resumes
		api.GET("/vacancies/upload", vacancyHandler.UploadVacancy)

		// Match routes
		api.POST("/match/:vacancy_id/match", matchHandler.MatchVacancyWithResumes)
		api.POST("/match/:vacancy_id/:resume_id", matchHandler.MatchVacancyWithResume)
		api.GET("/match/:vacancy_id/matches", matchHandler.GetVacancyResumeMatches)

		// More specific path to avoid conflict
		api.GET("/match/details/:match_id", matchHandler.GetVacancyResumeMatchByID)

		// One-time routes for updating embeddings (After DB switch)
		api.POST("/resumes/update-embeddings", resumeHandler.UpdateAllResumeEmbeddings)
		api.POST("/vacancies/update-embeddings", vacancyHandler.UpdateAllVacancyEmbeddings)
	}

	return r
}
