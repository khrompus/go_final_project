package main

import (
	"log"
	_ "modernc.org/sqlite"

	"github.com/khrompus/go_final_project/pkg/db"
	"github.com/khrompus/go_final_project/pkg/server"
)

func main() {

	storage, err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}()
	err = server.Run(storage)
	if err != nil {
		log.Fatal(err)
	}
}
