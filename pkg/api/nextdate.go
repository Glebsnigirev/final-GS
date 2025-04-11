package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

// NextDate вычисляет следующую дату на основе параметров
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим начальную дату
	parsedDStart, err := time.Parse(layout, dstart)
	if err != nil {
		return "", errors.New("invalid start date format")
	}

	// Разбираем правило повторения
	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		// Правило d <число>
		if len(parts) != 2 {
			return "", errors.New("incorrect 'd' rule format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("interval must be an integer")
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("interval must be between 1 and 400 days, got %d", interval)
		}
		for {
			nextDate := parsedDStart.AddDate(0, 0, interval)
			if nextDate.After(now) {
				return nextDate.Format(layout), nil
			}
			parsedDStart = nextDate
		}
	case "y":
		// Правило y: задача выполняется ежегодно
		for {
			nextDate := parsedDStart.AddDate(1, 0, 0) // Добавляем один год
			if nextDate.After(now) {                  // Проверяем, что дата больше текущей
				return nextDate.Format(layout), nil // Возврат даты в формате 20060102
			}
			parsedDStart = nextDate // Продолжаем искать следующую дату
		}
	case "w":
		// Реализуем еженедельные повторы
		weekdays := map[string]int{
			"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7,
		}
		validWeekdays := make([]int, 0)
		for _, wd := range strings.Split(parts[1], ",") {
			if weekday, ok := weekdays[wd]; ok {
				validWeekdays = append(validWeekdays, weekday)
			} else {
				return "", fmt.Errorf("invalid weekday value: %s", wd)
			}
		}
		for {
			nextDate := parsedDStart.AddDate(0, 0, 1)
			if nextDate.After(now) {
				for _, wd := range validWeekdays {
					if int(nextDate.Weekday()) == wd {
						return nextDate.Format(layout), nil
					}
				}
			}
			parsedDStart = nextDate
		}
	case "m":
		// Реализуем ежемесячные повторы
		daysAndMonths := strings.Split(parts[1], " ")
		if len(daysAndMonths) > 2 {
			return "", errors.New("too many arguments in 'm' rule")
		}
		daysStr := daysAndMonths[0]
		monthsStr := ""
		if len(daysAndMonths) == 2 {
			monthsStr = daysAndMonths[1]
		}
		validDays := make([]int, 0)
		for _, day := range strings.Split(daysStr, ",") {
			if dayInt, err := strconv.Atoi(day); err == nil {
				if dayInt >= 1 && dayInt <= 31 {
					validDays = append(validDays, dayInt)
				} else if dayInt == -1 || dayInt == -2 {
					validDays = append(validDays, dayInt)
				} else {
					return "", fmt.Errorf("invalid day value: %s", day)
				}
			} else {
				return "", fmt.Errorf("invalid day value: %s", day)
			}
		}
		validMonths := make([]int, 0)
		if monthsStr != "" {
			for _, month := range strings.Split(monthsStr, ",") {
				if monthInt, err := strconv.Atoi(month); err == nil {
					if monthInt >= 1 && monthInt <= 12 {
						validMonths = append(validMonths, monthInt)
					} else {
						return "", fmt.Errorf("invalid month value: %s", month)
					}
				} else {
					return "", fmt.Errorf("invalid month value: %s", month)
				}
			}
		} else {
			validMonths = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}
		for {
			nextDate := parsedDStart.AddDate(0, 0, 1)
			if nextDate.After(now) {
				if len(validMonths) > 0 && !contains(validMonths, int(nextDate.Month())) {
					continue
				}
				for _, day := range validDays {
					if day == -1 {
						lastDayOfMonth := time.Date(nextDate.Year(), nextDate.Month()+1, 0, 0, 0, 0, 0, time.UTC)
						if nextDate.Day() == lastDayOfMonth.Day() {
							return nextDate.Format(layout), nil
						}
					} else if day == -2 {
						lastDayOfPrevMonth := time.Date(nextDate.Year(), nextDate.Month(), 0, 0, 0, 0, 0, time.UTC)
						prevLastDay := lastDayOfPrevMonth.AddDate(0, 0, -1)
						if nextDate.Day() == prevLastDay.Day() {
							return nextDate.Format(layout), nil
						}
					} else if nextDate.Day() == day {
						return nextDate.Format(layout), nil
					}
				}
			}
			parsedDStart = nextDate
		}
	default:
		return "", errors.New("unsupported repeat rule type")
	}
}

// Вспомогательная функция для проверки наличия элемента в срезе
func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
