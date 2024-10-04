// internal/migration/migration.go

package migration

import (
	"clever_hr_api/internal/model"
	"log"

	"gorm.io/gorm"
)

func ApplyCustomMigrations(db *gorm.DB) error {
	// Define interview types with descriptions
	interviewTypes := []model.InterviewType{
		{Name: "Go", Description: "Interview for Go developers with focus on concurrency and design patterns."},
		{Name: "Python", Description: "Python interview focusing on web development, data analysis, and machine learning."},
		{Name: "Flutter", Description: "Mobile app development interview for Flutter developers with focus on UI/UX and state management."},
		{Name: "Java", Description: "Interview for Java developers with emphasis on enterprise solutions and multithreading."},
		{Name: "JavaScript", Description: "Front-end interview focusing on JavaScript, DOM manipulation, and frameworks like React."},
		{Name: "TypeScript", Description: "Interview for TypeScript developers with a focus on static typing and large-scale applications."},
		{Name: "Ruby on Rails", Description: "Web development interview for Ruby on Rails focusing on full-stack development."},
		{Name: "Kotlin", Description: "Interview for Kotlin developers with a focus on Android app development."},
		{Name: "C#", Description: "Interview focusing on C# for .NET developers, with emphasis on web apps and APIs."},
		{Name: "Swift", Description: "iOS app development interview for Swift developers with focus on SwiftUI and performance."},
		{Name: "PHP", Description: "Backend interview for PHP developers, focusing on server-side web apps and frameworks like Laravel."},
		{Name: "C++", Description: "Systems-level programming interview for C++ developers, focusing on memory management and performance."},
		{Name: "Rust", Description: "Rust interview for systems programming, emphasizing memory safety and concurrency."},
		{Name: "DevOps", Description: "DevOps interview focusing on CI/CD pipelines, cloud infrastructure, and automation."},
		{Name: "Machine Learning", Description: "Interview focusing on machine learning models, algorithms, and data preprocessing."},
		{Name: "Data Science", Description: "Data Science interview focusing on statistical analysis, data wrangling, and visualization."},
		{Name: "Frontend Development", Description: "Front-end interview focusing on HTML, CSS, JavaScript, and modern frameworks."},
		{Name: "Backend Development", Description: "Back-end development interview focusing on API design, databases, and security."},
		{Name: "Full Stack Development", Description: "Interview for full-stack developers with a focus on both front-end and back-end technologies."},
		{Name: "Cloud Computing", Description: "Interview focusing on cloud computing services (AWS, GCP, Azure) and scalable architectures."},
	}

	// Apply migrations
	for _, interviewType := range interviewTypes {
		if err := db.FirstOrCreate(&interviewType, model.InterviewType{Name: interviewType.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("Custom migrations applied successfully.")
	return nil
}
