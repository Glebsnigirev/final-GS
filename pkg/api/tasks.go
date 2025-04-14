package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db"
)

// taskHandler обрабатывает GET, POST и PUT запросы для задач
func taskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		// Получаем задачу по ID
		id := r.URL.Query().Get("id")
		if id == "" {
			writeErrorResponse(w, fmt.Errorf("Не указан идентификатор"))
			return
		}

		task, err := db.GetTask(id)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Ошибка: %v", err))
			return
		}

		// Для GET-запроса возвращаем объект типа *db.Task
		writeJson(w, task)

	case http.MethodPost:
		// Обрабатываем создание новой задачи
		var task db.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Ошибка при декодировании данных: %v", err))
			return
		}

		id, err := db.AddTask(&task)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Ошибка при добавлении задачи: %v", err))
			return
		}

		// Для POST-запроса возвращаем id добавленной задачи в виде map
		writeJson(w, map[string]interface{}{"id": id})

	case http.MethodPut:
		// Обновляем задачу
		id := r.URL.Query().Get("id")
		if id == "" {
			writeErrorResponse(w, fmt.Errorf("Не указан идентификатор"))
			return
		}

		var task db.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Ошибка при декодировании данных: %v", err))
			return
		}

		// Устанавливаем ID, которое мы получили из запроса
		task.ID = id

		// Логируем полученные данные для отладки
		log.Printf("Обновляем задачу с ID %s: %+v", id, task)

		// Проверяем, что дата не меньше текущей
		if task.Date < time.Now().Format("20060102") {
			writeErrorResponse(w, fmt.Errorf("Дата не может быть меньше сегодняшней"))
			return
		}

		// Обновляем задачу в базе данных
		err = db.UpdateTask(&task)
		if err != nil {
			writeErrorResponse(w, fmt.Errorf("Ошибка при обновлении задачи: %v", err))
			return
		}

		// Для PUT-запроса возвращаем пустой JSON
		writeJson(w, map[string]interface{}{})
	}
}

// writeJson пишет JSON в ответ
func writeJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
	}
}

// // writeErrorResponse пишет ошибку в формате JSON
// func writeErrorResponse(w http.ResponseWriter, err error) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusBadRequest)
// 	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
// }
