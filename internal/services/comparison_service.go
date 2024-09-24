package services

import (
	"1CHikNewProject/internal/database"
	"1CHikNewProject/internal/models"
	"math"
	"time"
)

func CompareAndSaveResults() error {
	var onesEvents []models.OnesEvent
	var hikEvents []models.HikvisionEvent

	// Получаем данные из таблиц
	database.DB.Find(&onesEvents)
	database.DB.Find(&hikEvents)

	for _, one := range onesEvents {
		for _, hik := range hikEvents {
			if one.TabNumber == hik.PersonID && one.Date == hik.EventTime.Truncate(24*time.Hour) {
				// Рассчитываем разницу
				difference := math.Abs(one.HoursWorked - hik.GetHours())

				// Сохраняем результат сравнения
				result := models.ComparisonResult{
					TabNumber:       one.TabNumber,
					EmployeeName:    one.EmployeeName,
					ObjectName:      one.ObjectName,
					Brigade:         one.Brigade,
					Date:            one.Date,
					OnesHoursWorked: one.HoursWorked,
					HikHoursWorked:  hik.GetHours(),
					Difference:      difference,
					NightShift:      one.NightShift,
				}
				database.DB.Create(&result)
			}
		}
	}
	return nil
}
