package api

import (
	"fmt"
	"strconv"
	"time"
)

// NextDate - вычисляет следующую дату с учетом правила повторения
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("Ошибка при парсинге даты: %v", err)
	}

	if repeat == "" {
		return taskDate.Format("20060102"), nil
	}

	if len(repeat) > 2 {
		switch repeat[:2] {
		case "d ":
			daysToAdd := repeat[2:]
			days, err := strconv.Atoi(daysToAdd)
			if err != nil {
				return "", fmt.Errorf("Неверный формат повторения (d N): %v", err)
			}
			nextDate := taskDate.AddDate(0, 0, days)
			return nextDate.Format("20060102"), nil
		case "w ":
			weeksToAdd := repeat[2:]
			weeks, err := strconv.Atoi(weeksToAdd)
			if err != nil {
				return "", fmt.Errorf("Неверный формат повторения (w N): %v", err)
			}
			nextDate := taskDate.AddDate(0, 0, weeks*7)
			return nextDate.Format("20060102"), nil
		case "m ":
			monthsToAdd := repeat[2:]
			months, err := strconv.Atoi(monthsToAdd)
			if err != nil {
				return "", fmt.Errorf("Неверный формат повторения (m N): %v", err)
			}
			nextDate := taskDate.AddDate(0, months, 0)
			return nextDate.Format("20060102"), nil
		case "y ":
			yearsToAdd := repeat[2:]
			years, err := strconv.Atoi(yearsToAdd)
			if err != nil {
				return "", fmt.Errorf("Неверный формат повторения (y N): %v", err)
			}
			nextDate := taskDate.AddDate(years, 0, 0)
			return nextDate.Format("20060102"), nil
		default:
			return "", fmt.Errorf("Невозможно обработать правило повторения: %v", repeat)
		}
	}

	return "", fmt.Errorf("Невозможно обработать правило повторения: %v", repeat)
}
