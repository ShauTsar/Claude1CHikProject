package main

import (
	"1CHikNewProject/internal/database"
	"1CHikNewProject/internal/handlers"
	"1CHikNewProject/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()
	//r.LoadHTMLGlob("templates/*")
	// Подключение HTML файлов
	r.LoadHTMLFiles("templates/correctTab.html", "templates/correctTabMinutes.html", "templates/detailTabV2.html")
	database.ConnectDB()

	// Маршруты для страниц
	r.GET("/version1", handlers.RenderVersion1)
	r.GET("/version1-minutes", handlers.RenderVersion1InMinutes)
	r.GET("/version2", handlers.RenderVersion2)
	// В cmd/main.go
	r.POST("/export", handlers.HandleExport)
	r.POST("/compare", handlers.CompareData)

	// Запуск сервера
	r.Run(":8889")
}
func StartDataSyncScheduler() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Запуск синхронизации данных")
				// Запуск функций для синхронизации
				//TODO добавить получения данных из таблиц
				services.FetchAndSave1CData("2024-09-01", "2024-09-30", []string{"Проект1", "Проект2"})
				services.FetchAndSaveHikVisionData("2024-09-01", "2024-09-30", []string{"door1", "door2"})
				services.CompareAndSaveResults()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
