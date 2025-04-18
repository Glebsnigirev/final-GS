package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()
	if err := checkDate(&task, now); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task, now time.Time) error {
	// Если дата указана как "today", устанавливаем текущую дату
	if task.Date == "today" {
		task.Date = now.Format("20060102")
	}

	// Если дата пустая, устанавливаем текущую дату
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	// Проверка корректности формата даты
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("Неверный формат даты")
	}

	// Если дата в прошлом, обновляем ее на текущую или на следующую по правилу повторения
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102") // Сегодняшняя дата
		} else {
			// Если есть правило повторения, вычисляем следующую дату
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("Неверное правило повторения: %v", err)
			}
			task.Date = next
		}
	}

	return nil
}
