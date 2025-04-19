package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const dateLayout = "20060102"

// Функция для проверки високосного года
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// Функция для вычисления следующей даты на основе правила повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка корректности повторения
	if err := validateRepeat(repeat); err != nil {
		return "", err
	}

	// Преобразование строки даты в тип time.Time
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	// Если дата уже прошла, добавляем годы или месяцы
	if startDate.Before(now) {
		if repeat == "y" {
			// Добавляем год, пока дата не станет после текущей
			for startDate.Before(now) {
				startDate = startDate.AddDate(1, 0, 0)
			}
		} else if repeat == "m" {
			// Добавляем месяц, пока дата не станет после текущей
			for startDate.Before(now) {
				startDate = startDate.AddDate(0, 1, 0)
			}
		}
	}

	// Если правило повторения по дням (например, "d 10")
	re := regexp.MustCompile(`^d\s*(\d+)$`)
	match := re.FindStringSubmatch(repeat)
	if len(match) > 0 {
		daysToAdd, err := strconv.Atoi(match[1])
		if err != nil {
			return "", fmt.Errorf("invalid day count")
		}

		// Добавляем дни и проверяем, чтобы дата не прошла
		startDate = startDate.AddDate(0, 0, daysToAdd)

		// Проверяем, что дата не прошла
		for startDate.Before(now) {
			startDate = startDate.AddDate(0, 0, daysToAdd)
		}

		return startDate.Format("20060102"), nil
	}

	// Проверка на високосный год
	if startDate.Month() == time.February && startDate.Day() == 29 && !isLeapYear(startDate.Year()) {
		startDate = time.Date(startDate.Year(), time.March, 1, 0, 0, 0, 0, startDate.Location()) // Переносим на 1 марта
	}

	// Возвращаем дату в формате YYYYMMDD
	return startDate.Format("20060102"), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	// Обработка даты "now"
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	// Вычисление следующей даты
	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправка ответа
	fmt.Fprint(w, next)
}
