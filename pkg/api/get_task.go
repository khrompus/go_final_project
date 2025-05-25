package api

import (
	"net/http"
)

func (dBase *API) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := dBase.storage.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}

	writeJson(w, task)
}
