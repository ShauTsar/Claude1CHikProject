package handlers

import (
	"1CHikNewProject/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CompareData сравнивает данные из 1С и HikVision
func CompareData(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	projects := c.QueryArray("projects")

	result, err := services.Compare1CAndHik(startDate, endDate, projects)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleData обрабатывает запрос на сравнение данных между 1С и HikVision
func HandleData(c *gin.Context) {
	var form struct {
		StartDate        string   `form:"startDate"`
		EndDate          string   `form:"endDate"`
		SelectedProjects []string `form:"projects"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные формы"})
		return
	}

	fmt.Printf("Получены данные: месяц=%s, проекты=%v\n", form.StartDate, form.SelectedProjects)

	usersFrom1C, err := services.Search1C(form.StartDate, form.EndDate, form.SelectedProjects)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users1C := services.ProcessUsers1C(usersFrom1C)
	doorIndexCodes := services.GetDoorIndexCodesForProjects(form.SelectedProjects)

	hikEvents, err := services.GetAllHikVisionEvents(form.StartDate, form.EndDate, doorIndexCodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usersHik := services.ProcessUsersHik(hikEvents, form.SelectedProjects)

	response := struct {
		Users1C  map[string]services.User1C  `json:"users1C"`
		UsersHik map[string]services.UserHik `json:"usersHik"`
	}{
		Users1C:  users1C,
		UsersHik: usersHik,
	}

	c.JSON(http.StatusOK, response)
}
