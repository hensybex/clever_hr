// migration/migration.go

package migration

import (
	"clever_hr_api/internal/model"
	"clever_hr_api/internal/model/categories"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ApplyCustomMigrations(db *gorm.DB) error {
	// Define default user credentials (you should hash the password)
	defaultUsername := "admin"
	defaultPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Define the default user
	defaultUser := model.User{
		Username: defaultUsername,
		Password: string(hashedPassword), // Store hashed password
	}

	// Apply user migration (check if user already exists by username)
	if err := db.FirstOrCreate(&defaultUser, model.User{Username: defaultUsername}).Error; err != nil {
		return err
	}

	groups := []category_model.JobGroup{
		{
			ID:      1,
			TitleRu: "Разработка программного обеспечения",
			TitleEn: "Software Development",
		},
		{
			ID:      2,
			TitleRu: "Контроль качества, тестирование",
			TitleEn: "Quality Assurance",
		},
		{
			ID:      3,
			TitleRu: "Аналитика данных",
			TitleEn: "Data Analytics",
		},
		{
			ID:      4,
			TitleRu: "Управление проектами",
			TitleEn: "Project Management",
		},
	}

	specializations := []category_model.Specialization{
		{
			ID:         1,
			TitleRu:    "Бэкенд разработчик",
			TitleEn:    "Backend Web Developer",
			JobGroupID: 1,
		},
		{
			ID:         2,
			TitleRu:    "Фронтенд разработчик",
			TitleEn:    "Frontend Web Developer",
			JobGroupID: 1,
		},
		{
			ID:         3,
			TitleRu:    "Мобильный разработчик",
			TitleEn:    "Mobile Developer",
			JobGroupID: 1,
		},
		{
			ID:         4,
			TitleRu:    "DevOps инженер",
			TitleEn:    "DevOps Engineer",
			JobGroupID: 1,
		},
		{
			ID:         5,
			TitleRu:    "Инженер по автомат. тестированию",
			TitleEn:    "Test Automation Engineer",
			JobGroupID: 2,
		},
		{
			ID:         6,
			TitleRu:    "Инженер по производительности",
			TitleEn:    "Software Performance Engineer",
			JobGroupID: 2,
		},
		{
			ID:         7,
			TitleRu:    "Ручной тестировщик",
			TitleEn:    "Manual Tester",
			JobGroupID: 2,
		},
		{
			ID:         8,
			TitleRu:    "Дата-аналитик",
			TitleEn:    "Data Analyst",
			JobGroupID: 3,
		},
		{
			ID:         9,
			TitleRu:    "Дата-сайентист",
			TitleEn:    "Data Scientist",
			JobGroupID: 3,
		},
		{
			ID:         10,
			TitleRu:    "Инженер по данным",
			TitleEn:    "Data Engineer",
			JobGroupID: 3,
		},
		{
			ID:         11,
			TitleRu:    "Учёный-исследователь",
			TitleEn:    "Research Scientist",
			JobGroupID: 3,
		},
		{
			ID:         12,
			TitleRu:    "Менеджер проектов",
			TitleEn:    "Project Manager",
			JobGroupID: 4,
		},
		{
			ID:         13,
			TitleRu:    "Менеджер по продуктам",
			TitleEn:    "Product Manager",
			JobGroupID: 4,
		},
		{
			ID:         14,
			TitleRu:    "ТимЛид",
			TitleEn:    "Team Lead",
			JobGroupID: 4,
		},
	}

	// Predefined Qualifications
	qualifications := []category_model.Qualification{
		{
			ID:      1,
			TitleRu: "Стажёр",
			TitleEn: "Intern",
		},
		{
			ID:      2,
			TitleRu: "Ассистент",
			TitleEn: "Assistant",
		},
		{
			ID:      3,
			TitleRu: "Младший специалист",
			TitleEn: "Junior Specialist",
		},
		{
			ID:      4,
			TitleRu: "Специалист",
			TitleEn: "Specialist",
		},
		{
			ID:      5,
			TitleRu: "Старший специалист",
			TitleEn: "Senior Specialist",
		},
		{
			ID:      6,
			TitleRu: "Ведущий специалист",
			TitleEn: "Lead Specialist",
		},
		{
			ID:      7,
			TitleRu: "Директор",
			TitleEn: "Director",
		},
	}

	for _, group := range groups {
		if err := db.FirstOrCreate(&group, category_model.JobGroup{ID: group.ID}).Error; err != nil {
			return err
		}
	}

	for _, specialization := range specializations {
		if err := db.FirstOrCreate(&specialization, category_model.Specialization{ID: specialization.ID}).Error; err != nil {
			return err
		}
	}

	for _, qualification := range qualifications {
		if err := db.FirstOrCreate(&qualification, category_model.Qualification{ID: qualification.ID}).Error; err != nil {
			return err
		}
	}

	log.Println("Custom migrations applied successfully.")
	return nil
}
