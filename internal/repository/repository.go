package repository

import (
	"github.com/dkshi/hephaestus"
	"github.com/jmoiron/sqlx"
)

type Task interface {
	CreateTask(task *hephaestus.Task) (int64, error)
	CreateTaskTable() error
	GetTasks(chatId int64) ([]hephaestus.Task, error)
	DeleteTask(taskId int64) error
}

type Repository struct {
	Task
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Task: NewTaskPostgres(db),
	}
}

func (r *Repository) CreateTables() error {
	if err := r.CreateTaskTable(); err != nil {
		return err
	}
	return nil
}
