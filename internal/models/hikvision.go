package models

import (
	"fmt"
	"gorm.io/gorm"
	"sort"
	"time"
)

// EventDetails хранит данные о событиях (время и URL фото)
type EventDetails struct {
	Time   time.Time `json:"time"`
	PicUri string    `json:"pic_uri"`
}

type Event struct {
	EventID           string `json:"eventId"`
	EventType         int    `json:"eventType"`
	EventTime         string `json:"eventTime"`
	PersonID          string `json:"personId"`
	PersonName        string `json:"personName"`
	PersonType        string `json:"personType"`
	DoorName          string `json:"doorName"`
	DoorIndexCode     string `json:"doorIndexCode"`
	CardNo            string `json:"cardNo"`
	CheckInAndOutType int    `json:"checkInAndOutType"`
	PicUri            string `json:"picUri"`
	DeviceTime        string `json:"deviceTime"`
	TemperatureData   string `json:"temperatureData"`
	TemperatureStatus int    `json:"temperatureStatus"`
	WearMaskStatus    int    `json:"wearMaskStatus"`
	ReaderIndexCode   string `json:"readerIndexCode"`
	ReaderName        string `json:"readerName"`
}

// UserHik представляет данные пользователя из HikVision
type UserHik struct {
	PersonID   string                  `json:"person_id"`   // Уникальный ID пользователя
	PersonName string                  `json:"person_name"` // Имя пользователя
	Dates      map[string]interface{}  `json:"dates"`       // Время работы по дням (может быть строка, часы или другие данные)
	FirstEntry map[string]EventDetails `json:"first_entry"` // Первое событие (вход)
	LastExit   map[string]EventDetails `json:"last_exit"`   // Последнее событие (выход)
	Object     string                  `json:"object"`      // Объект (проект)
}

// HikvisionEvent представляет событие из HikVision
type HikvisionEvent struct {
	gorm.Model
	PersonID          string    `json:"person_id"`
	PersonName        string    `json:"person_name"`
	DoorName          string    `json:"door_name"`
	DoorIndexCode     string    `json:"door_index_code"`
	EventType         int       `json:"event_type"`
	EventTime         time.Time `json:"event_time"`
	CheckInAndOutType int       `json:"check_in_and_out_type"`
	PicUri            string    `json:"pic_uri"`
	NightShift        bool      `json:"night_shift"`
}

// CalculateWorkHours рассчитывает отработанные часы для пользователя HikVision
func CalculateWorkHours(user *UserHik, date string, events []time.Time) {
	if len(events) == 0 {
		return
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Before(events[j])
	})

	var workStart, workEnd time.Time
	if len(events) >= 2 {
		for _, t := range events {
			if t.Hour() >= 4 {
				workStart = t
				break
			}
		}
		if workStart.IsZero() {
			workStart = events[0]
		}
		workEnd = events[len(events)-1]
		if workEnd.Before(workStart) {
			workEnd = workEnd.AddDate(0, 0, 1)
		}

		duration := workEnd.Sub(workStart)
		if duration > 16*time.Hour {
			duration = 16 * time.Hour
		}

		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		user.Dates[date] = fmt.Sprintf("%02d:%02d", hours, minutes)
	} else {
		t := events[0]
		if t.Hour() < 12 {
			user.Dates[date] = fmt.Sprintf("П %02d:%02d", t.Hour(), t.Minute())
		} else {
			user.Dates[date] = fmt.Sprintf("У %02d:%02d", t.Hour(), t.Minute())
		}
	}
}
