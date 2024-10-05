// internal/router/router.go

package router

import (
	"clever_hr_api/internal/handlers"
	"clever_hr_api/internal/repository"
	"clever_hr_api/internal/service"
	"clever_hr_api/internal/usecase"

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

	// Initialize services
	mistralService := service.NewMistralService(gptCallRepo)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	candidateUsecase := usecase.NewCandidateUsecase(candidateRepo, resumeRepo, resumeAnalysisResultRepo)
	resumeUsecase := usecase.NewResumeUsecase(resumeRepo, resumeAnalysisResultRepo, candidateRepo, userRepo, *mistralService)
	interviewUsecase := usecase.NewInterviewUsecase(interviewRepo, interviewTypeRepo, interviewMessageRepo, interviewAnalysisResultRepo, resumeRepo, *mistralService)
	interviewTypeUsecase := usecase.NewInterviewTypeUsecase(interviewTypeRepo)
	resumeAnalysisResultUsecase := usecase.NewResumeAnalysisResultUsecase(resumeAnalysisResultRepo, candidateRepo, resumeRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUsecase)
	resumeHandler := handlers.NewResumeHandler(resumeUsecase, userUsecase, resumeAnalysisResultUsecase)
	interviewHandler := handlers.NewInterviewHandler(interviewUsecase)
	candidateHandler := handlers.NewCandidateHandler(candidateUsecase)
	interviewTypeHandler := handlers.NewInterviewTypeHandler(interviewTypeUsecase)

	// Define API routes
	api := r.Group("/api")
	{
		// WebSocket routes for interview analysis
		api.GET("/ws/interview/analyse", interviewHandler.AnalyseInterviewMessageWebsocket)

		// User routes
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:user_id/info", userHandler.GetUser)
		api.PUT("/users/:user_id/switch", userHandler.SwitchUserType)
		api.GET("/users/:user_id/candidates", userHandler.GetCandidatesByUserID)
		api.GET("/users/:user_id/get_role", userHandler.GetUserRoleByTgID)

		// Resume routes
		api.POST("/resumes/employee", resumeHandler.EmployeeUploadResume)
		api.POST("/resumes/candidate", resumeHandler.CandidateUploadResume)
		api.GET("/resumes/:resume_id/analyze", resumeHandler.RunResumeAnalysis)
		api.GET("/resumes/:resume_id/analysis-result", resumeHandler.GetResumeAnalysisResult)

		// Interview routes
		api.POST("/interviews", interviewHandler.CreateInterview)
		api.GET("/interviews/:interview_id/analyze", interviewHandler.RunFullInterviewAnalysis)
		api.GET("/interviews/:interview_id/analysis-result", interviewHandler.GetInterviewAnalysisResult)

		// Candidates routes
		api.GET("/candidates/:candidate_id", candidateHandler.GetCandidateInfo)
		api.GET("/candidates/:candidate_id/resume", candidateHandler.GetResume)

		// Interview type routes
		api.GET("/interview-types", interviewTypeHandler.ListInterviewTypes)
	}

	return r
}
