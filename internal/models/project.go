package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name string `json:"name"`
}

type Door struct {
	gorm.Model
	DoorIndexCode string `json:"door_index_code"`
	ProjectID     uint   `json:"project_id"`
}
