package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db" // Импорт db
)

// Init регистрирует маршруты для API
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/add", addTaskHandler) // Добавляем обработчик для добавления задачи
}

// nextDayHandler обрабатывает запросы к маршруту /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	// Если параметр now не определен, берем текущую дату
	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowParam) // Здесь уточнил формат
		if err != nil {
			http.Error(w, "Invalid 'now' parameter", http.StatusBadRequest)
			return
		}
	}

	// Парсим дату
	_, err := time.Parse("20060102", dateParam) // Уточнил формат
	if err != nil {
		http.Error(w, "Invalid 'date' parameter", http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату
	// Передаем now как time.Time объект, а не строку
	nextDate, err := db.NextDate(now, dateParam, repeatParam) // Передаем time.Time
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат в формате 20060102
	fmt.Fprintf(w, "%s", nextDate)
}
