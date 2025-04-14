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

	// Проверка формата повторений
	if len(repeat) < 3 {
		return fmt.Errorf("Невозможное повторение: %v", repeat)
	}

	switch repeat[:2] {
	case "d ", "w ", "m ", "y ":
		// Пытаемся распарсить число после префикса
		_, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return fmt.Errorf("Неверный формат повторения: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("Невозможное повторение: %v", repeat)
	}
}

func isValidDate(date string) (bool, error) {
	// Пробуем распарсить строку в дату с использованием формата "20060102"
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		// Логируем ошибку и возвращаем ошибку с подробностями
		return false, fmt.Errorf("неправильный формат даты: %v", err)
	}

	// Получаем текущую дату без времени (00:00:00)
	currentDate := time.Now().Truncate(24 * time.Hour)

	// Логируем текущую дату и дату задачи для отладки
	log.Printf("Проверка даты: задача - %v, текущая - %v", parsedDate.Format("20060102"), currentDate.Format("20060102"))

	// Проверяем, что дата не меньше текущей
	if parsedDate.Before(currentDate) {
		// Логируем ошибку с датами и возвращаем ошибку
		return false, fmt.Errorf("Дата %v не может быть меньше сегодняшней (%v)", date, currentDate.Format("20060102"))
	}

	// Если дата валидна, возвращаем true
	return true, nil
}

// AddTask добавляет новую задачу в базу данных
func AddTask(task *Task) (string, error) {
	// Проверка на пустое значение для даты
	if task.Date == "" {
		task.Date = time.Now().Format("20060102") // Устанавливаем текущую дату, если она пустая
	}

	// Проверка даты
	valid, err := isValidDate(task.Date)
	if !valid {
		log.Printf("Ошибка при проверке даты: %v", err) // Логируем ошибку
		return "", err
	}

	// Проверка повторений
	err = validateRepeat(task.Repeat)
	if err != nil {
		log.Printf("Ошибка при проверке повторения: %v", err) // Логируем ошибку
		return "", err
	}

	// Вставляем задачу в базу данных
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи в БД: %v", err) // Логируем ошибку
		return "", err
	}

	// Получаем ID добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Ошибка при получении ID задачи: %v", err) // Логируем ошибку
		return "", err
	}

	// Возвращаем ID задачи в виде строки
	log.Printf("Задача добавлена с ID: %d", id) // Логируем успешное добавление
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
