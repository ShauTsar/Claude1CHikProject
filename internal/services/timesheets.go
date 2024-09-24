package services

import (
	"1CHikNewProject/internal/database"
	"1CHikNewProject/internal/models"
	"log"
)

// CreateTimesheet creates a new timesheet record
func CreateTimesheet(timesheet *models.Timesheet) error {
	result := database.DB.Create(timesheet)
	if result.Error != nil {
		log.Printf("Error creating timesheet: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetTimesheets retrieves all timesheets
func GetTimesheets() ([]models.Timesheet, error) {
	var timesheets []models.Timesheet
	result := database.DB.Find(&timesheets)
	if result.Error != nil {
		log.Printf("Error fetching timesheets: %v", result.Error)
		return nil, result.Error
	}
	return timesheets, nil
}
