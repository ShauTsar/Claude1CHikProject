package services

import (
	"1CHikNewProject/internal/sdk"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	tabsFrom1CUrl = "http://10.150.0.6/es-hrm-kadr/hs/tabs/get_tabs/"
)

type PersonInfo struct {
	PersonCode string `json:"person_code"`
	PersonName string `json:"person_name"`
	PersonID   string `json:"person_id"`
}

type Response struct {
	Data struct {
		PersonCode string `json:"person_code"`
		PersonName string `json:"person_name"`
		PersonId   string `json:"person_id"`
	} `json:"data"`
}

const getUserPersonIdByTab = "api/path/to/get/person/info"

// GetPersonInfoByCode извлекает данные о пользователе из HikVision по коду сотрудника
func GetPersonInfoByCode(personCode string) (*PersonInfo, error) {
	body := map[string]interface{}{
		"personCode": personCode,
	}

	result, err := workWithHikAPI(body, getUserPersonIdByTab)
	if err != nil {
		return nil, err
	}

	resJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return nil, err
	}

	var response Response
	err = json.Unmarshal(resJson, &response)
	if err != nil {
		return nil, err
	}
	return &PersonInfo{PersonCode: response.Data.PersonCode, PersonName: response.Data.PersonName, PersonID: response.Data.PersonId}, nil
}

// Search1C ищет данные в 1С за указанный период и по выбранным проектам
func Search1C(startDate, endDate string, allowedProjects []string) ([]User, error) {
	return search1C(startDate, endDate, allowedProjects)
}

type User struct {
	Date     string `json:"date"`
	Object   string `json:"object"`
	Brigade  string `json:"brigade"`
	Employee string `json:"employee"`
	Tab      string `json:"tab"`
	Time     string `json:"time"`
	Hours    int    `json:"hours"`
}

func search1C(startDate, endDate string, allowedProjects []string) ([]User, error) {
	fmt.Printf("search1C вызвана с параметрами: месяц=%s,%s, проекты=%v\n", startDate, endDate, allowedProjects)

	if startDate == "" || len(allowedProjects) == 0 {
		return nil, fmt.Errorf("некорректные параметры: месяц или проекты не указаны")
	}
	url := tabsFrom1CUrl + startDate + "/" + endDate
	// Кодирование логина и пароля в формате base64
	auth := "НовиковНА:02nikita"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	// Создание HTTP-запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return nil, err
	}
	// Добавление заголовка аутентификации
	req.Header.Add("Authorization", "Basic "+encodedAuth)
	// Отправка HTTP-запроса
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("неожиданный статус ответа: %d, тело: %s", resp.StatusCode, string(body))
	}
	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return nil, err
	}
	var allUsers []User
	err = json.Unmarshal(body, &allUsers)
	if err != nil {
		fmt.Println("Ошибка разбора JSON:", err)
		return nil, err
	}

	projectMap := make(map[string]bool)
	for _, project := range allowedProjects {
		projectMap[project] = true
	}

	var filteredUsers []User
	for _, user := range allUsers {
		if projectMap[user.Object] {
			filteredUsers = append(filteredUsers, user)
		}
	}
	if len(filteredUsers) == 0 {
		fmt.Println("Не найдено пользователей, соответствующих критериям")
	}

	fmt.Printf("search1C возвращает %d пользователей\n", len(filteredUsers))
	return filteredUsers, nil
}

func workWithHikAPI(body map[string]interface{}, apiURL string) (sdk.Result, error) {
	hk := sdk.HKConfig{
		Ip:      "10.150.0.15",
		Port:    10443,
		AppKey:  "29086835",
		Secret:  "Eoah9YBzDTYIVW67gRCL",
		IsHttps: true,
	}
	result, err := hk.HttpPost(apiURL, body, 15)
	if err != nil {
		return result, err
	}
	return result, nil
}
func workWithHikAPIPictures(body map[string]interface{}, apiURL string) (string, error) {
	hk := sdk.HKConfig{
		Ip:      "10.150.0.15",
		Port:    10443,
		AppKey:  "29086835",
		Secret:  "Eoah9YBzDTYIVW67gRCL",
		IsHttps: true,
	}
	result, err := hk.HttpPostPictures(apiURL, body, 20)
	if err != nil {
		return result, err
	}
	return result, nil
}
