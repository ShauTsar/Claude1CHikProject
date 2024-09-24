package handlers

import (
	"1CHikNewProject/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RenderVersion1 рендерит первую версию страницы сравнения
func RenderVersion1(c *gin.Context) {
	c.HTML(http.StatusOK, "correctTab.html", nil)
}

// RenderVersion1InMinutes рендерит первую версию страницы в минутах
func RenderVersion1InMinutes(c *gin.Context) {
	c.HTML(http.StatusOK, "correctTabMinutes.html", nil)
}

func RenderVersion2(c *gin.Context) {
	projects, err := services.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить проекты"})
		return
	}

	c.HTML(http.StatusOK, "detailTabV2.html", gin.H{
		"Projects": projects,
	})
}
