package api

import (
	"encoding/json"
	"fmt"
	"github.com/khrompus/go_final_project/pkg/db"
	"net/http"
	"strings"
	"time"
)

func afterNow(now, t time.Time) bool {
	// Нормализуем даты
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return nowDate.After(tDate)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	// Проверяем пустая ли строка с датой
	if len(task.Date) == 0 {
		task.Date = now.Format("20060102")
		return nil
	}
	// Проверяем валидность даты
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты. Ожидается YYYYMMDD")
	}

	//Сравниваем время сейчас и task.Date
	// Если правило повтора не указано, просто проверяем дату
	if task.Repeat == "" {
		if now.After(t) {
			task.Date = now.Format("20060102")
		}
		return nil
	}
	//Вычисляем следующую дату по правилу Repeat
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("неверное правило повтора: %v", err)
	}
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, берём сегодняшнее число
			task.Date = now.Format("20060102")
		} else {
			// иначе берём вычисленную следующую дату
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON-объект с ID вместо просто числа
	writeJson(w, struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}
