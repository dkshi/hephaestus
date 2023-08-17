package hephaestus

import "time"

type Task struct {
	TaskID int64
	ChatID   int64
	TaskName string
	Deadline time.Time
}
