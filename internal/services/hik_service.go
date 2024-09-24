package services

import (
	"1CHikNewProject/internal/database"
	"1CHikNewProject/internal/models"
	"encoding/json"
	"time"
)

func processUser(personID string, events []models.Event, layout string, doorProjectMap map[string]string) models.UserHik {
	user := models.UserHik{
		PersonID:   personID,
		PersonName: "",
		Dates:      make(map[string]interface{}),
		FirstEntry: make(map[string]models.EventDetails),
		LastExit:   make(map[string]models.EventDetails),
		Object:     "",
	}

	dayEvents := make(map[string][]time.Time)

	for _, event := range events {
		eventTime, _ := time.Parse(layout, event.EventTime)
		date := eventTime.Format("2006-01-02")

		user.PersonName = event.PersonName

		if firstEntry, exists := user.FirstEntry[date]; !exists || eventTime.Before(firstEntry.Time) {
			user.FirstEntry[date] = models.EventDetails{Time: eventTime, PicUri: event.PicUri}
		}
		if lastExit, exists := user.LastExit[date]; !exists || eventTime.After(lastExit.Time) {
			user.LastExit[date] = models.EventDetails{Time: eventTime, PicUri: event.PicUri}
		}
		if project, exists := doorProjectMap[event.DoorIndexCode]; exists {
			user.Object = project
		}

		dayEvents[date] = append(dayEvents[date], eventTime)
	}

	for date, events := range dayEvents {
		models.CalculateWorkHours(&user, date, events)
	}

	return user
}
func FetchAndSaveHikVisionData(startDate, endDate string, doorIndexCodes []string) error {
	eventTypes := []int{196893, 197161, 197151}

	for _, eventType := range eventTypes {
		pageNo := 1
		pageSize := 500

		for {
			body := map[string]interface{}{
				"startTime":      startDate + "T00:00:00+08:00",
				"endTime":        endDate + "T23:59:59+08:00",
				"eventType":      eventType,
				"doorIndexCodes": doorIndexCodes,
				"pageNo":         pageNo,
				"pageSize":       pageSize,
			}

			result, err := workWithHikAPI(body, getEventsByPersonURL)
			if err != nil {
				return err
			}

			resJson, err := json.Marshal(result)
			if err != nil {
				return err
			}

			var response models.ResponseDataPersonEvents
			if err := json.Unmarshal(resJson, &response); err != nil {
				return err
			}

			for _, event := range response.Data.List {
				eventTime, _ := time.Parse("2006-01-02T15:04:05", event.EventTime)

				// Сохраняем данные в таблицу hikvision_events
				hikEvent := models.HikvisionEvent{
					PersonID:          event.PersonID,
					PersonName:        event.PersonName,
					DoorName:          event.DoorName,
					DoorIndexCode:     event.DoorIndexCode,
					EventType:         event.EventType,
					EventTime:         eventTime,
					CheckInAndOutType: event.CheckInAndOutType,
					PicUri:            event.PicUri,
					NightShift:        isNightShift(event.EventTime),
				}
				database.DB.Create(&hikEvent)
			}

			if len(response.Data.List) < pageSize {
				break
			}
			pageNo++
		}
	}
	return nil
}
