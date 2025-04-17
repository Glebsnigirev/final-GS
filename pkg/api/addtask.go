package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db" // Импорт db
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Дата представлена в формате, отличном от 20060102"})
		return
	}

	// Если дата в прошлом:
	if t.Before(now.Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			writeJson(w, map[string]string{"error": "Дата не может быть меньше сегодняшней"})
			return
		}

		// Есть повтор — пересчитываем дату
		nextDate, err := db.NextDate(t, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка при вычислении следующей даты: %v", err)})
			return
		}
		task.Date = nextDate
	}

	// Добавляем задачу
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка при добавлении задачи: %v", err)})
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// Вспомогательная функция для записи JSON в ответ
func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Ошибка записи JSON", http.StatusInternalServerError)
	}
}
