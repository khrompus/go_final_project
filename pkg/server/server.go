package server

import (
	"fmt"
	"github.com/khrompus/go_final_project/pkg/api"
	"net/http"
)

func Run() error {
	port := 7540

	//запускаем api
	r := api.Init()
	r.Handle("/*", http.FileServer(http.Dir("web")))
	fmt.Println("Приложение успешно запущено")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
