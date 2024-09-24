package models

import (
	"gorm.io/gorm"
)

// ComparisonResult хранит результат сравнения 1С и HikVision
type ComparisonResult struct {
	gorm.Model
	HikvisionEventID uint    `json:"hikvision_event_id"`
	OnesEventID      uint    `json:"ones_event_id"`
	ComparisonDate   string  `json:"comparison_date"`
	IsMatched        bool    `json:"is_matched"`
	HikvisionHours   float64 `json:"hikvision_hours"`
	OnesHours        float64 `json:"ones_hours"`
	PhotoURL         string  `json:"photo_url"`
}
