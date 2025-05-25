package server

import (
	"fmt"
	"net/http"

	"github.com/khrompus/go_final_project/pkg/api"
	"github.com/khrompus/go_final_project/pkg/db"
)

func Run(storage *db.TaskStorage) error {
	port := 7540

	api := api.NewAPI(storage)
	//запускаем api
	r := api.Init()
	r.Mount("/", http.FileServer(http.Dir("web")))
	fmt.Println("Приложение успешно запущено")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
