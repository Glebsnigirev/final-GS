package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Убираем время из now
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	startDate, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty repeat rule")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid daily format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		for startDate.Before(now) {
			startDate = startDate.AddDate(0, 0, days)
		}
		return startDate.Format(dateLayout), nil
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid yearly format")
		}
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(now) {
				break
			}
		}
		return startDate.Format(dateLayout), nil
	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
