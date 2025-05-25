package db

import (
	_ "database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Tasks получаем задачи по заданному лимиту
func (s *TaskStorage) Tasks(limit int) ([]Task, error) {

	query := `SELECT id, date, title, comment, repeat 
              FROM scheduler 
              ORDER BY date ASC 
              LIMIT ?`
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return tasks, fmt.Errorf("failed to scan task row: %v", err)
		}
		task.ID = fmt.Sprintf("%d", id) // Преобразуем ID в строку
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return []Task{}, fmt.Errorf("error after iterating rows: %v", err)
	}
	return tasks, nil
}

// AddTask добавляем задачу в бд
func (s *TaskStorage) AddTask(task *Task) (int64, error) {
	var id int64
	// Запрос в бд на создание задачи
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// GetTask Получаем задачу по id
func (s *TaskStorage) GetTask(id string) (*Task, error) {
	query := `SELECT * FROM scheduler WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var task Task

	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return nil, fmt.Errorf("failed to scan task row: %v", err)
	}

	return &task, nil
}

// UpdateTask Обновляем данные задачи
func (s *TaskStorage) UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
              SET date = ?, title = ?, comment = ?, repeat = ? 
              WHERE id = ?`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (s *TaskStorage) DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (s *TaskStorage) UpdateDate(next string, id string) error {
	query := `UPDATE scheduler 
              SET date = ? WHERE id = ?`
	res, err := s.db.Exec(query, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
