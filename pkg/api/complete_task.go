package api

import (
	"net/http"
	"time"
)

func (dBase *API) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}
	task, err := dBase.storage.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusInternalServerError)
	}
	// Если Repeat пустой, то удаляем задачу
	if task.Repeat == "" {
		if err := dBase.storage.DeleteTask(id); err != nil {
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
		if err := dBase.storage.UpdateDate(next, id); err != nil {
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
