package main

import (
	"fmt"
	"github.com/khrompus/go_final_project/pkg/db"
	"github.com/khrompus/go_final_project/pkg/server"
	"log"
)

func main() {
	fmt.Println("Hello world")

	err := server.Run()
	if err != nil {
		panic(err)
	}
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.CloseDb(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	// Далее можно работать с базой через GetDB()
	fmt.Println("Приложение успешно запущено")
}
