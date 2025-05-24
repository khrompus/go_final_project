package api

import (
	"github.com/khrompus/go_final_project/pkg/db"
	"net/http"
	"time"
)

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusInternalServerError)
	}
	// Если Repeat пустой, то удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			if err.Error() == "task not found" {
				writeError(w, "Задача не найдена", http.StatusNotFound)
			} else {
				writeError(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		writeJson(w, struct{}{})
	} else {
		//Если Repeat не пустой, обновляем дату
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка изменения даты", http.StatusInternalServerError)
		}
		if err := db.UpdateDate(next, id); err != nil {
			if err.Error() == "task not found" {
				writeError(w, "Задача не найдена", http.StatusNotFound)
			} else {
				writeError(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		writeJson(w, struct{}{})
	}
}
