package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX idx_date ON scheduler(date);
`

// Init открывает соединение с базой данных и создает схему при необходимости
func Init(dbFile string) error {
	// Проверяем существование файла
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %v", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе: %v", err)
	}

	// Создаем таблицы при первом запуске
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %v", err)
		}
		fmt.Println("База данных инициализирована (созданы таблицы)")
	}

	return nil
}
func CloseDb() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDB возвращает глобальное соединение с базой данных
func GetDB() *sql.DB {
	return db
}
