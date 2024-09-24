package services

import (
	"1CHikNewProject/internal/database"
	"1CHikNewProject/internal/models"
)

// GetProjects возвращает список всех проектов
func GetProjects() ([]models.Project, error) {
	var projects []models.Project
	result := database.DB.Find(&projects)
	return projects, result.Error
}

// GetProjectsFromDB возвращает список проектов из базы данных
func GetProjectsFromDB() ([]models.Project, error) {
	var projects []models.Project
	result := database.DB.Find(&projects)
	return projects, result.Error
}

// GetDoorsFromDB возвращает список дверей из базы данных
func GetDoorsFromDB() ([]models.Door, error) {
	var doors []models.Door
	result := database.DB.Find(&doors)
	return doors, result.Error
}
