package api

import (
	"fmt"
	"github.com/khrompus/go_final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"` // Используем структуру из пакета db
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, fmt.Sprintf("failed to get tasks: %v", err), http.StatusInternalServerError)
		return
	}

	// Конвертируем []db.Task в []*db.Task
	taskPointers := make([]*db.Task, len(tasks))
	for i := range tasks {
		taskPointers[i] = &tasks[i]
	}

	writeJson(w, TasksResp{
		Tasks: taskPointers,
	})
}
