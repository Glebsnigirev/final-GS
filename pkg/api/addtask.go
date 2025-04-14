package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db"
)

const layoutISO = "20060102"

// Функция для возврата ошибок в формате JSON
func writeErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// Обработчик POST-запросов для добавления задач
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение и логирование тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("Ошибка чтения тела запроса: %v", err))
		log.Printf("Ошибка чтения тела запроса: %v", err)
		return
	}
	log.Printf("Тело запроса: %s", string(body)) // Логируем запрос

	// Десериализация JSON из запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorResponse(w, fmt.Errorf("Ошибка десериализации JSON: %v", err))
		log.Printf("Ошибка десериализации JSON: %v", err) // Логируем ошибку
		return
	}
	log.Printf("Задача после десериализации: %+v", task) // Логируем десериализованный объект

	// Далее обработка задачи как обычно...

	// Проверка обязательных полей
	if task.Title == "" {
		writeErrorResponse(w, fmt.Errorf("Не указан заголовок задачи"))
		return
	}

	// Проверка формата даты
	now := time.Now()
	var t time.Time
	if task.Date == "" {
		task.Date = now.Format(layoutISO)
		t = now
	} else {
		var err error
		t, err = time.Parse(layoutISO, task.Date)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Неверный формат даты: %s", task.Date))
			return
		}
	}

	// Если дата задачи меньше текущей, обновляем на сегодняшнюю, но только если не задано повторение
	if t.Before(now) && task.Repeat == "" {
		task.Date = now.Format(layoutISO)
	} else if t.Before(now) && task.Repeat != "" {
		// Если задано повторение, оставляем дату как есть и рассчитываем следующее повторение
		var nextDate string
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Ошибка при вычислении даты повторения: %v", err))
			log.Printf("Ошибка при вычислении следующей даты для задачи: %v", err) // Логируем ошибку
			return
		}
		nextDate = next
		task.Date = nextDate
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("Ошибка добавления задачи в БД: %v", err))
		log.Printf("Ошибка при добавлении задачи в БД: %v", err) // Логируем ошибку
		return
	}

	// Ответ с идентификатором задачи
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}
