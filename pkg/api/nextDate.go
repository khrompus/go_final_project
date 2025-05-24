package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	// Получаем параметры из запроса
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeatRule := query.Get("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(res, fmt.Sprintf("invalid now parameter: %v", err), http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(res, "date parameter is required", http.StatusBadRequest)
		return
	}
	if repeatRule == "" {
		http.Error(res, "repeat parameter is required", http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(now, dateStr, repeatRule)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(res, nextDate)
}

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	const layout = "20060102"
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
