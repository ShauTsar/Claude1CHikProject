package tests

import (
	"1CHikNewProject/internal/models"
	"1CHikNewProject/internal/services"
	"testing"
)

func TestCreateTimesheet(t *testing.T) {
	timesheet := models.Timesheet{
		EmployeeID:  "12345",
		Date:        "2023-09-23",
		HoursWorked: 8.0,
		Task:        "Development",
	}

	err := services.CreateTimesheet(&timesheet)
	if err != nil {
		t.Fatalf("Failed to create timesheet: %v", err)
	}
}

func TestGetTimesheets(t *testing.T) {
	timesheets, err := services.GetTimesheets()
	if err != nil {
		t.Fatalf("Failed to fetch timesheets: %v", err)
	}

	if len(timesheets) == 0 {
		t.Fatalf("Expected timesheets, got none")
	}
}
