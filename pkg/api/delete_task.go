package api

import (
	"net/http"
)

func (dBase *API) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//  Достаем ID
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}
	//	Удаляем задачу
	if err := dBase.storage.DeleteTask(id); err != nil {
		if err.Error() == "task not found" {
			writeError(w, "Задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем пустой JSON объект
	writeJson(w, struct{}{})
}
