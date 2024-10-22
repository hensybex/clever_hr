// cmd/server/main.go

package main

import (
	"clever_hr_api/internal/config"
	"clever_hr_api/internal/migration"
	"clever_hr_api/internal/model"
	category_model "clever_hr_api/internal/model/categories"
	"clever_hr_api/internal/router"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	milvusClient := config.InitMilvusClient()
	defer milvusClient.Close()

	if err := model.CreateMilvusCollections(milvusClient); err != nil {
		log.Fatalf("Failed to create Milvus collections: %v", err)
	}

	// AutoMigrate models
	if err := db.AutoMigrate(
		&model.GPTCall{},
		&model.User{},
		&model.Resume{},
		//&model.ResumeEmbedding{},
		&model.Vacancy{},
		//&model.VacancyEmbedding{},
		&model.ResumeAnalysisResult{},
		&model.VacancyResumeMatch{},
		&category_model.JobGroup{},
		&category_model.Qualification{},
		&category_model.Specialization{},
		&model.Message{},
	); err != nil {
		log.Fatalf("Failed to automigrate: %v", err)
	}

	err = migration.ApplyCustomMigrations(db)
	if err != nil {
		log.Fatalf("Failed to apply custom migration : %v", err)
	}

	// Setup router
	r := router.SetupRouter(db, milvusClient)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
