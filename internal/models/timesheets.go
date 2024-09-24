package models

import (
	"gorm.io/gorm"
)

// Timesheet represents a timesheet record
type Timesheet struct {
	gorm.Model
	EmployeeID  string  `json:"employee_id"`
	Date        string  `json:"date"`
	HoursWorked float64 `json:"hours_worked"`
	Task        string  `json:"task"`
}
