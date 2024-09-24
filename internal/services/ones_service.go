package services

import (
	"1CHikNewProject/internal/database"
	"1CHikNewProject/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func FetchAndSave1CData(startDate, endDate string, allowedProjects []string) error {
	url := fmt.Sprintf("%s%s/%s", tabsFrom1CUrl, startDate, endDate)
	auth := "НовиковНА:02nikita"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var users []models.User1C
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return err
	}

	// Фильтрация по проектам
	projectMap := make(map[string]bool)
	for _, project := range allowedProjects {
		projectMap[project] = true
	}

	for _, user := range users {
		if projectMap[user.Object] {
			date, _ := time.Parse("2006-01-02T15:04:05", user.Date)

			// Определяем ночную смену
			nightShift := isNightShift(user.Time)

			// Сохраняем данные в таблицу ones_events
			event := models.OnesEvent{
				TabNumber:    user.Tab,
				EmployeeName: user.Employee,
				ObjectName:   user.Object,
				Brigade:      user.Brigade,
				Date:         date,
				HoursWorked:  float64(user.Hours),
				TimeType:     user.Time,
				NightShift:   nightShift,
				ProjectID:    GetProjectIDByObjectName(user.Object), // Получаем ID проекта
			}
			database.DB.Create(&event)
		}
	}
	return nil
}

func isNightShift(timeType string) bool {
	// Логика определения ночной смены
	return timeType == "НС"
}
