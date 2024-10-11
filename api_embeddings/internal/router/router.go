// internal/router/router.go

package router

import (
	"clever_hr_embeddings/internal/handlers"
	"clever_hr_embeddings/internal/repository"
	"clever_hr_embeddings/internal/usecases"

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
	embeddingsRepo := repository.NewEmbeddingRepository(db)
	resumeRepo := repository.NewResumeRepository(db)
	vacancyRepo := repository.NewVacancyRepository(db)

	/* // Initialize services
	mistralService := mistral.NewMistralService(gptCallRepo) */

	// Initialize use cases
	candidateMatchingUsecase := usecases.NewCandidateMatchingUsecase(resumeRepo, vacancyRepo, embeddingsRepo)
	//resumeAnalysisUsecase := usecases.NewResumeAnalysisUsecase()
	resumeUsecase := usecases.NewResumeUsecase(resumeRepo, embeddingsRepo)
	vacancyUsecase := usecases.NewVacancyUsecase(vacancyRepo, embeddingsRepo)
	//vacancyAnalysisUsecase := usecases.NewCandidateMatchingUsecase(resumeRepo, vacancyRepo, embeddingsRepo)

	// Define API routes
	api := r.Group("/api")
	{
		// WebSocket routes for interview analysis
		api.POST("/resumes", handlers.UploadResumeHandler(resumeUsecase))
		api.POST("/vacancies", handlers.PostVacancyHandler(vacancyUsecase))
		api.GET("/vacancies/:id/candidates", handlers.GetBestCandidatesHandler(candidateMatchingUsecase))
	}

	return r
}
