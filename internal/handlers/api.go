package handlers

import (
	"1CHikNewProject/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetHikVisionPersonInfo получает информацию о пользователе из HikVision по его коду
func GetHikVisionPersonInfo(c *gin.Context) {
	personCode := c.Param("personCode")
	personInfo, err := services.GetPersonInfoByCode(personCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, personInfo)
}

// Get1CData получает данные из 1С по диапазону дат и проектам
func Get1CData(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	allowedProjects := c.QueryArray("projects")

	if startDate == "" || endDate == "" || len(allowedProjects) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	users, err := services.Search1C(startDate, endDate, allowedProjects)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
