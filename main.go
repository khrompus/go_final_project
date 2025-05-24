package main

import (
	"fmt"
	"github.com/khrompus/go_final_project/pkg/db"
	"github.com/khrompus/go_final_project/pkg/server"
	"log"
	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("Hello world")

	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	err := server.Run()
	if err != nil {
		panic(err)
	}

	// Далее можно работать с базой через GetDB()
}
