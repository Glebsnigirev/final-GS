package api

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/glebsnigirev/final-GS/pkg/db"
// )

// func taskHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		addTaskHandler(w, r)
// 	default:
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 	}
// }

// func addTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	var task db.Task

// 	// Десериализуем JSON из тела запроса
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&task)
// 	if err != nil {
// 		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
// 		return
// 	}

// 	// Проверяем обязательное поле title
// 	if task.Title == "" {
// 		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
// 		return
// 	}

// 	// Получаем текущую дату
// 	now := time.Now()

// 	// Проверяем и подставляем текущую дату, если поле date пустое
// 	if task.Date == "" {
// 		task.Date = now.Format("20060102")
// 	}

// 	// Проверка формата даты
// 	t, err := time.Parse("20060102", task.Date)
// 	if err != nil {
// 		writeJson(w, map[string]string{"error": "Дата представлена в формате, отличном от 20060102"})
// 		return
// 	}

// 	// Если задача в прошлом, корректируем её дату
// 	var next string
// 	if afterNow(now, t) {
// 		// Если указано правило повторения, вычисляем следующую дату
// 		if task.Repeat != "" {
// 			next, err = NextDate(now, task.Date, task.Repeat)
// 			if err != nil {
// 				writeJson(w, map[string]string{"error": "Неверный формат повторения"})
// 				return
// 			}
// 			task.Date = next
// 		} else {
// 			// Если правило повторения не указано, подставляем сегодняшнюю дату
// 			task.Date = now.Format("20060102")
// 		}
// 	}

// 	// Добавляем задачу в базу данных
// 	id, err := db.AddTask(&task)
// 	if err != nil {
// 		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка при добавлении задачи: %v", err)})
// 		return
// 	}

// 	// Возвращаем JSON с ID добавленной задачи
// 	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
// }

// // Вспомогательная функция для записи JSON в ответ
// func writeJson(w http.ResponseWriter, data interface{}) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	encoder := json.NewEncoder(w)
// 	err := encoder.Encode(data)
// 	if err != nil {
// 		http.Error(w, "Ошибка записи JSON", http.StatusInternalServerError)
// 	}
// }
