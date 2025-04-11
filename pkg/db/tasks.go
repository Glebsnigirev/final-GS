package db

import (
	"fmt"
	"regexp"
	"time"
)

func validateDate(date string) error {
	// Проверка формата даты
	match, err := regexp.MatchString(`^\d{8}$`, date)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("invalid date format: expected YYYYMMDD, got %s", date)
	}

	// Проверка корректности даты
	_, err = time.Parse("20060102", date)
	if err != nil {
		return fmt.Errorf("invalid date: %s", err)
	}

	return nil
}

func validateTask(task *Task) error {
	// Проверка даты
	if err := validateDate(task.Date); err != nil {
		return err
	}

	// Проверка обязательного поля Title
	if task.Title == "" {
		return fmt.Errorf("Title cannot be empty")
	}

	// Проверка максимальной длины Title
	if len(task.Title) > 255 {
		return fmt.Errorf("Title exceeds maximum length of 255 characters")
	}

	// Проверка максимальной длины Comment
	if len(task.Comment) > 1000 {
		return fmt.Errorf("Comment exceeds maximum length of 1000 characters")
	}

	return nil
}

func AddTask(task *Task) (int64, error) {
	// Проверка данных
	if err := validateTask(task); err != nil {
		return 0, err
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %v", err)
	}

	return id, nil
}
