// cmd/server/main.go

package main

import (
	"clever_hr_api/internal/config"
	"clever_hr_api/internal/migration"
	"clever_hr_api/internal/model"
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

	// AutoMigrate models
	if err := db.AutoMigrate(
		&model.GPTCall{},
		&model.User{},
		&model.Candidate{},
		&model.Resume{},
		&model.ResumeAnalysisResult{},
		&model.InterviewType{},
		&model.Interview{},
		&model.InterviewMessage{},
		&model.InterviewAnalysisResult{},
	); err != nil {
		log.Fatalf("Failed to automigrate: %v", err)
	}

	err = migration.ApplyCustomMigrations(db)
	if err != nil {
		log.Fatalf("Failed to apply custom migration : %v", err)
	}

	// Setup router
	r := router.SetupRouter(db)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
