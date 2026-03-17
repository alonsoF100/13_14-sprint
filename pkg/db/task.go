package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, errors.New("база данных не инициализирована")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetTasks(limit int) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("база данных не инициализирована")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, errors.New("база данных не инициализирована")
	}

	var taskID int64
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("неверный формат идентификатора")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	task := &Task{}
	var dbID int64

	err = DB.QueryRow(query, taskID).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}

	task.ID = strconv.FormatInt(dbID, 10)
	return task, nil
}

func UpdateTask(task *Task) error {
	if DB == nil {
		return errors.New("база данных не инициализирована")
	}

	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("неверный формат идентификатора")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, taskID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с id %s не найдена", task.ID)
	}

	return nil
}
