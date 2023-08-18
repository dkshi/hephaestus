package hephaestus

import "time"

type Task struct {
	TaskID   int64     `db:"task_id"`
	ChatID   int64     `db:"chat_id"`
	TaskName string    `db:"task_name"`
	Deadline time.Time `db:"deadline"`
}
