package repository

import (
	"errors"
	"fmt"

	"github.com/dkshi/hephaestus"
	"github.com/jmoiron/sqlx"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (t *TaskPostgres) CreateTask(task *hephaestus.Task) (int64, error) {
	var taskId int64

	query := fmt.Sprintf("INSERT INTO %s (chat_id, task_name, deadline) values ($1, $2, $3) RETURNING task_id;", tasksTable)
	row := t.db.QueryRow(query, task.ChatID, task.TaskName, task.Deadline)

	if err := row.Scan(&taskId); err != nil {
		return 0, err
	}

	return taskId, nil
}

func (t *TaskPostgres) CreateTaskTable() error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (task_id INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY, chat_id VARCHAR(256) NOT NULL, task_name VARCHAR(256) NOT NULL, deadline TIMESTAMP NOT NULL);", tasksTable)
	if _, err := t.db.Exec(query); err != nil {
		return err
	}
	return nil
}

func (t *TaskPostgres) GetTasks(chatId int64) ([]hephaestus.Task, error) {
	var tasks []hephaestus.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE chat_id = $1", tasksTable)

	if err := t.db.Select(&tasks, query, chatId); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TaskPostgres) DeleteTask(taskId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE task_id = $1", tasksTable)
	res, err := t.db.Exec(query, taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("error rows affected: rowsAffected = 0")
	}
	if err != nil {
		return err
	}

	return nil
}
