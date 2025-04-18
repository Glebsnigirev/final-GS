package api

import (
	"log"
	"net/http"
)

func Init() {
	// Регистрируем только нужные обработчики
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request:", r.Method, r.URL)

	switch r.Method {
	case http.MethodPost:
		log.Println("Handling POST request")
		addTaskHandler(w, r)
	}
}
