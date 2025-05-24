package api

import (
	"encoding/json"
	"github.com/khrompus/go_final_project/pkg/db"
	"net/http"
	"strings"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Неправильный формат JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(w, "Заголовок не может быть пустым", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if err.Error() == "task not found" {
			writeError(w, "Задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJson(w, struct{}{})
}
