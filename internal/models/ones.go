package models

import (
	"gorm.io/gorm"
	"time"
)

// OnesEvent представляет событие из 1С
type OnesEvent struct {
	gorm.Model
	TabNumber    string    `json:"tab_number"`
	EmployeeName string    `json:"employee_name"`
	ObjectName   string    `json:"object_name"`
	Brigade      string    `json:"brigade"`
	Date         time.Time `json:"date"`
	HoursWorked  float64   `json:"hours_worked"`
	TimeType     string    `json:"time_type"`
	NightShift   bool      `json:"night_shift"`
	ProjectID    uint      `json:"project_id"`
}
type User1C struct {
	Date     string `json:"date"`     // Дата события
	Object   string `json:"object"`   // Объект (проект)
	Brigade  string `json:"brigade"`  // Бригада
	Employee string `json:"employee"` // ФИО сотрудника
	Tab      string `json:"tab"`      // Табельный номер
	Time     string `json:"time"`     // Тип смены (например, "НС")
	Hours    int    `json:"hours"`    // Количество отработанных часов
}
