package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	parsedDStart, err := time.Parse(layout, dstart)
	if err != nil {
		return "", errors.New("invalid start date format")
	}

	// Разделяем правило на части
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 || parts[0] == "" {
		return "", errors.New("empty repeat rule")
	}

	ruleType := parts[0]
	switch ruleType {
	case "d":
		// Правило d <число>
		if len(parts) != 2 {
			return "", errors.New("incorrect 'd' rule format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("interval must be an integer")
		}

		// Проверка максимального интервала перемещена сюда, чтобы она выполнялась ДО начала цикла
		if interval > 400 {
			return "", errors.New("maximum allowed interval is 400 days")
		}

		// Исправлено: не возвращаем ошибку, если интервал равен 1 и дата совпадает с текущей
		if interval == 1 && parsedDStart.Format(layout) == now.Format(layout) {
			return parsedDStart.Format(layout), nil
		}

		// Начало основного цикла для ежедневного увеличения даты
		for {
			nextDate := parsedDStart.AddDate(0, 0, interval)
			if nextDate.After(now) {
				return nextDate.Format(layout), nil
			}
			parsedDStart = nextDate
		}
	case "y":
		// Правило y
		for {
			nextDate := parsedDStart.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate.Format(layout), nil
			}
			parsedDStart = nextDate
		}
	default:
		return "", errors.New("unsupported repeat rule type")
	}
}
