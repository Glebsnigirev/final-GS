package api

import (
	"fmt"
	"net/http"
	"time"
)

// Init регистрирует маршруты для API
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
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
		now, err = time.Parse(layout, nowParam)
		if err != nil {
			http.Error(w, "Invalid 'now' parameter", http.StatusBadRequest)
			return
		}
	}

	// Парсим дату
	_, err := time.Parse(layout, dateParam)
	if err != nil {
		http.Error(w, "Invalid 'date' parameter", http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат в формате 20060102
	fmt.Fprintf(w, "%s", nextDate)
}
