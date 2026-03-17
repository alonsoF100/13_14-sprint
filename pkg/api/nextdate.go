package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	start, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", date)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		return handleDaily(now, start, parts)
	case "y":
		return handleYearly(now, start)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", parts[0])
	}
}

func handleDaily(now time.Time, start time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("не указан интервал для правила d")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("неверный формат интервала: %s", parts[1])
	}

	if interval < 1 || interval > 400 {
		return "", fmt.Errorf("интервал должен быть от 1 до 400 дней")
	}

	next := start
	for {
		next = next.AddDate(0, 0, interval)
		if afterNow(next, now) {
			break
		}
		if next.Year() > now.Year()+10 {
			return "", fmt.Errorf("не удалось найти следующую дату")
		}
	}

	return next.Format(DateFormat), nil
}

func handleYearly(now time.Time, start time.Time) (string, error) {
	next := start
	for {
		next = next.AddDate(1, 0, 0)
		if afterNow(next, now) {
			break
		}
		if next.Year() > now.Year()+100 {
			return "", fmt.Errorf("не удалось найти следующую дату")
		}
	}
	return next.Format(DateFormat), nil
}

func afterNow(date time.Time, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	if y1 > y2 {
		return true
	}
	if y1 < y2 {
		return false
	}
	if m1 > m2 {
		return true
	}
	if m1 < m2 {
		return false
	}
	return d1 > d2
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	if dateParam == "" {
		http.Error(w, "параметр date обязателен", http.StatusBadRequest)
		return
	}
	if repeatParam == "" {
		http.Error(w, "параметр repeat обязателен", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
