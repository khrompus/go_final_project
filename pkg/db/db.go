package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type TaskStorage struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX idx_date ON scheduler(date);
`

// Init открывает соединение с базой данных и создает схему при необходимости
func Init(dbFile string) (*TaskStorage, error) {
	// Проверяем существование файла
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	fmt.Println("Файл существует:", !install)
	// Открываем базу данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %v", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе: %v", err)
	}

	// Создаем таблицы при первом запуске
	if install {
		if _, err := db.Exec(schema); err != nil {
			return nil, fmt.Errorf("ошибка создания схемы: %v", err)
		}
		fmt.Println("База данных инициализирована (созданы таблицы)")
	}

	return &TaskStorage{db: db}, nil
}
func (s *TaskStorage) Close() error {
	return s.db.Close()
}
