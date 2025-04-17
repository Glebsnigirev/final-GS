package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate - вычисляет следующую дату с учетом правила повторения
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	const maxDays = 400

	// Парсим дату из строки
	taskDate, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", errors.New("ошибка при парсинге даты")
	}

	// Если правило повторения пустое, возвращаем пустую строку
	if repeat == "" {
		return "", nil
	}

	// Обрабатываем правило "y" отдельно
	if repeat == "y" {
		nextDate := taskDate.AddDate(1, 0, 0)
		if !afterNow(nextDate, now) {
			// Продолжаем добавлять годы, пока дата не окажется в будущем
			for !afterNow(nextDate, now) {
				nextDate = nextDate.AddDate(1, 0, 0)
			}
		}
		return nextDate.Format("20060102"), nil
	}

	// Другие правила обрабатываем через Split
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("некорректный формат правила повторения: %v", repeat)
	}

	command := parts[0]
	value := parts[1]

	// Анализируем правило повторения
	switch command {
	case "d":
		days, err := strconv.Atoi(value)
		if err != nil || days > maxDays {
			return "", fmt.Errorf("неверный формат повторения (d N): %v", err)
		}

		// Сдвигаем как минимум на один интервал
		nextDate := taskDate.AddDate(0, 0, days)
		for !afterNow(nextDate, now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", fmt.Errorf("невозможно обработать правило повторения: %v", repeat)
	}
}

// afterNow проверяет, больше ли дата date, чем now, игнорируя время суток
func afterNow(date time.Time, now time.Time) bool {
	// Приводим обе даты к состоянию без часов, минут и секунд
	cleanedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	cleanedNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Сравниваем очищенные даты
	return cleanedDate.After(cleanedNow) || cleanedDate.Equal(cleanedNow)
}
