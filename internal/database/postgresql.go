package database

import (
	"1CHikNewProject/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// ConnectDB initializes the PostgreSQL database connection
func ConnectDB() {
	dsn := "host=localhost user=postgres password=NNA2s*123 dbname=hik_1c_timesheets port=15432 sslmode=disable TimeZone=UTC"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Автоматически мигрируем все модели
	err = DB.AutoMigrate(&models.UserHik{}, &models.User1C{}, &models.ComparisonResult{}, &models.Project{}, &models.Project{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database connection established and migrations applied")
}
