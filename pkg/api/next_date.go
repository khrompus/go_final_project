package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()

	// Получаем параметры из запроса
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeatRule := query.Get("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(layout, nowStr)
		if err != nil {
			writeError(w, fmt.Sprintf("invalid now parameter: %v", err), http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		writeError(w, "date parameter is required", http.StatusBadRequest)
		return
	}
	if repeatRule == "" {
		writeError(w, "repeat parameter is required", http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(now, dateStr, repeatRule)
	if err != nil {
		// Определяем тип ошибки для выбора HTTP статуса
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "invalid date format") ||
			strings.Contains(err.Error(), "invalid repeat rule") {
			status = http.StatusUnprocessableEntity
		}
		writeError(w, err.Error(), status)
		return
	}
	fmt.Fprint(w, nextDate)
}

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	//Проверка на наличие Repeat
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	// Нормализуем now (убираем время кроме даты)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	//Парсинг Даты
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}
	// Нормализуем date (убираем время кроме даты)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	switch {
	case repeat == "y":
		// Годовое повторение
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format(layout), nil

	case strings.HasPrefix(repeat, "d "):
		// Дневное повторение
		daysStr := strings.TrimPrefix(repeat, "d ")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid days interval")
		}

		// Для "d 1" просто добавляем 1 день пока не получим будущую дату
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}
		return date.Format(layout), nil

	default:
		return "", errors.New("invalid repeat rule")
	}
}
