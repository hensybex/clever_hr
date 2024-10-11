// internal/router/router.go

package router

import (
	"clever_hr_api/internal/handlers"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/service"
	"clever_hr_api/internal/usecase"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
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
	candidateRepo := repository.NewCandidateRepository(db)
	resumeRepo := repository.NewResumeRepository(db)
	resumeAnalysisResultRepo := repository.NewResumeAnalysisResultRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)
	interviewMessageRepo := repository.NewInterviewMessageRepository(db)
	interviewAnalysisResultRepo := repository.NewInterviewAnalysisResultRepository(db)
	gptCallRepo := repository.NewGPTCallRepository(db)
	interviewTypeRepo := repository.NewInterviewTypeRepository(db)
	vacancyRepo := repository.NewVacancyRepository(db)
	embeddingRepo := repository.NewEmbeddingRepository(db)

	// Initialize services
	mistralService := service.NewMistralService(gptCallRepo)
	authService := service.NewAuthService(userRepo, nil)

	// Initialize AuthMiddleware
	authMiddleware, err := AuthMiddleware(authService)
	if err != nil {
		log.Fatal("JWT Error: ", err)
	}

	// Set the jwtMiddleware in AuthService after it's created
	authService = service.NewAuthService(userRepo, authMiddleware)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	candidateUsecase := usecase.NewCandidateUsecase(candidateRepo, resumeRepo, userRepo, resumeAnalysisResultRepo, interviewRepo, interviewAnalysisResultRepo)
	resumeUsecase := usecase.NewResumeUsecase(resumeRepo, resumeAnalysisResultRepo, candidateRepo, userRepo, *mistralService)
	interviewUsecase := usecase.NewInterviewUsecase(interviewRepo, interviewTypeRepo, interviewMessageRepo, interviewAnalysisResultRepo, resumeRepo, userRepo, candidateRepo, *mistralService)
	interviewTypeUsecase := usecase.NewInterviewTypeUsecase(interviewTypeRepo)
	resumeAnalysisResultUsecase := usecase.NewResumeAnalysisResultUsecase(resumeAnalysisResultRepo, candidateRepo, resumeRepo)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo, embeddingRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUsecase, candidateUsecase)
	resumeHandler := handlers.NewResumeHandler(resumeUsecase, userUsecase, resumeAnalysisResultUsecase)
	interviewHandler := handlers.NewInterviewHandler(interviewUsecase)
	candidateHandler := handlers.NewCandidateHandler(candidateUsecase)
	interviewTypeHandler := handlers.NewInterviewTypeHandler(interviewTypeUsecase)
	vacancyHandler := handlers.NewVacancyHandler(vacancyUsecase)

	// Public routes (no auth needed)
	r.POST("/login", authMiddleware.LoginHandler)

	// Define API routes
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		// WebSocket routes for interview analysis
		api.GET("/ws/interview/analyse", interviewHandler.AnalyseInterviewMessageWebsocket)

		// User routes
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:user_id/info", userHandler.GetUser)
		api.PUT("/users/:user_id/switch", userHandler.SwitchUserType)
		//api.GET("/users/:user_id/candidates", userHandler.GetCandidatesByUserID)
		//api.GET("/users/:user_id/get_role", userHandler.GetUserRoleByTgID)

		// Resume routes
		api.POST("/resumes/upload", resumeHandler.UploadResume)
		//api.POST("/resumes", resumeHandler.CreateResumeHandler)
		api.GET("/resumes/:resume_id/analyze", resumeHandler.RunResumeAnalysis)
		api.GET("/resumes/:resume_id/analysis-result", resumeHandler.GetResumeAnalysisResult)

		// Vacancy routes
		api.POST("/vacancies", vacancyHandler.PostVacancy)
		api.GET("/vacancies", vacancyHandler.GetVacancies)
		api.GET("/vacancies/:id", vacancyHandler.GetVacancyByID)
		api.PUT("/vacancies/:id/status", vacancyHandler.UpdateVacancyStatus)

		// Interview routes
		api.POST("/interviews", interviewHandler.CreateInterview)
		api.GET("/interviews/:interview_id/analyze", interviewHandler.RunFullInterviewAnalysis)
		api.GET("/interviews/:interview_id/analysis-result", interviewHandler.GetInterviewAnalysisResult)

		// Candidates routes
		api.GET("/candidates/:candidate_id/get_by_id", candidateHandler.GetCandidateInfo)
		api.GET("/candidates/:candidate_id/resume", candidateHandler.GetResume)

		// Interview type routes
		api.GET("/interview-types", interviewTypeHandler.ListInterviewTypes)
	}

	return r
}
