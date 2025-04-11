package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/glebsnigirev/go_final_project_GS/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Читаем тело запроса
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	// Проверяем обязательные поля
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверяем дату
	if err := validateDate(task.Date); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Проверяем и применяем правило повторения
	if task.Repeat != "" {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		task.Date = nextDate
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возврат успешного результата
	writeJSON(w, map[string]interface{}{"id": strconv.FormatInt(id, 10)})
}

// Вспомогательная функция для отправки JSON-ответа
func writeJSON(w http.ResponseWriter, data interface{}) {
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

// Функция для валидации даты
func validateDate(date string) error {
	// Проверка формата даты
	match, err := regexp.MatchString(`^\d{8}$`, date)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("invalid date format: expected YYYYMMDD, got %s", date)
	}

	// Проверка корректности даты
	_, err = time.Parse("20060102", date)
	if err != nil {
		return fmt.Errorf("invalid date: %s", err)
	}

	return nil
}
