package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/khrompus/go_final_project/pkg/db"
)

type API struct {
	storage *db.TaskStorage
}

func NewAPI(storage *db.TaskStorage) *API {
	return &API{storage: storage}
}

func (api *API) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/nextdate", nextDayHandler)
	r.Post("/api/task", api.addTaskHandler)
	r.Get("/api/tasks", api.tasksHandler)
	r.Get("/api/task", api.getTaskHandler)
	r.Put("/api/task", api.updateTaskHandler)
	r.Delete("/api/task", api.deleteTaskHandler)
	r.Post("/api/task/done", api.completeTaskHandler)
	return r
}
