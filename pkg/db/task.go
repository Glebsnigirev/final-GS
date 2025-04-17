package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Task представляет собой структуру для хранения информации о задаче
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func validateRepeat(repeat string) error {
	if repeat == "" {
		return nil
	}

	// Проверка формата повторений: префикс и число
	if len(repeat) < 3 {
		return fmt.Errorf("Невозможное повторение: %v", repeat)
	}

	switch repeat[:2] {
	case "d ", "w ", "m ", "y ":
		if len(repeat) < 3 {
			return fmt.Errorf("Невозможное повторение: %v", repeat)
		}
		_, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return fmt.Errorf("Неверный формат повторения: %v", repeat)
		}
		return nil
	default:
		return fmt.Errorf("Невозможное повторение: %v", repeat)
	}
}

// AddTask добавляет новую задачу в базу данных
func AddTask(task *Task) (string, error) {
	// Устанавливаем текущую дату, если не указана
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// Проверка даты на корректность
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Printf("Неверный формат даты: %v", err)
		return "", fmt.Errorf("Неверный формат даты: %v", task.Date)
	}

	// Проверка повторений
	err = validateRepeat(task.Repeat)
	if err != nil {
		log.Printf("Ошибка при проверке повторения: %v", err)
		return "", err
	}

	// Вставляем задачу в базу данных
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи в БД: %v", err)
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		// Логируем ошибку, если не можем получить ID
		log.Printf("Ошибка при получении ID задачи: %v", err)

		// Попробуем проверить другой способ получения ID (например, RowsAffected)
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Ошибка при получении количества затронутых строк: %v", err)
			return "", err
		}
		if rowsAffected == 0 {
			return "", fmt.Errorf("Не удалось вставить задачу в базу данных")
		}
		return "", fmt.Errorf("Не удалось получить ID задачи: %v", err)
	}

	log.Printf("Задача добавлена с ID: %d", id)
	return fmt.Sprintf("%d", id), nil
}

// GetTask получает задачу по ID
func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, err
	}
	return &task, nil
}

// UpdateTask обновляет задачу в базе данных
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена для обновления")
	}

	return nil
}
