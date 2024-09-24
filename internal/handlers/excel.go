package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// handleExport экспортирует данные в Excel
func HandleExport(c *gin.Context) {
	var tableData struct {
		Headers []string   `json:"headers"`
		Data    [][]string `json:"data"`
	}

	if err := c.ShouldBindJSON(&tableData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при разборе данных"})
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Сравнение"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания листа в Excel"})
		return
	}
	f.SetActiveSheet(index)

	// Запись заголовков
	for col, header := range tableData.Headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		cell := colName + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// Запись данных
	for row, rowData := range tableData.Data {
		for col, cellData := range rowData {
			colName, _ := excelize.ColumnNumberToName(col + 1)
			cell := colName + string(row+2)
			f.SetCellValue(sheetName, cell, cellData)

			// Подсветка различий, если разница больше 30 минут
			if strings.Contains(tableData.Headers[col], "Hik") && col > 0 {
				hours1C, err1 := parseHours(rowData[col-1])
				hoursHik, err2 := parseHours(cellData)
				if err1 == nil && err2 == nil {
					diff := math.Abs(hours1C - hoursHik)
					if diff > 0.5 { // Разница больше 30 минут
						style, _ := f.NewStyle(&excelize.Style{
							Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
						})
						f.SetCellStyle(sheetName, cell, cell, style)
					}
				}
			}
		}
	}

	// Установка стилей для заголовков
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	lastColName, _ := excelize.ColumnNumberToName(len(tableData.Headers))
	f.SetCellStyle(sheetName, "A1", lastColName+"1", headerStyle)

	// Автоподбор ширины столбцов
	for col := 1; col <= len(tableData.Headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Отправка файла клиенту
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=compare.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи файла"})
		return
	}
}

func parseHours(input string) (float64, error) {
	if input == "" {
		return 0, nil
	}
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("некорректный формат времени")
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return float64(hours) + float64(minutes)/60, nil
}
