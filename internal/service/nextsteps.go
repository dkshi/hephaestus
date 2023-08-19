package service

import (
	"fmt"
	"strconv"

	"github.com/dkshi/hephaestus"
	"github.com/dkshi/hephaestus/scripts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Creating linked list with steps needed to procede (only for command /new)
func (s *Service) linkedCmdCreate(chatId int64) {
	s.LastStep = &nextStepLinkedList{
		val: s.CommandCreateSecond,
		next: &nextStepLinkedList{
			val:  s.CommandCreateLast,
			next: nil,
		},
	}

	s.task = &hephaestus.Task{
		ChatID: chatId,
	}
}

func (s *Service) linkedCmdComplete() {
	s.LastStep = &nextStepLinkedList{
		val:  s.CommandCompleteLast,
		next: nil,
	}
}

// Handles next steps for command /new
func (s *Service) CommandCreateSecond(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) {
	m.Text = "Enter deadline for your new task"

	taskName := upd.Message.Text
	s.task.TaskName = taskName

	s.LastStep = s.LastStep.next
}

func (s *Service) CommandCreateLast(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) {
	taskDL := upd.Message.Text
	parsedTaskDL, err := scripts.ParseStringToTime(taskDL)
	m.Text = "Incorrect date format, try: YYYY-MM-DD HH:MM:SS"

	if err != nil {
		return
	}

	s.task.Deadline = parsedTaskDL
	taskId, err := s.repo.CreateTask(s.task)

	if err != nil {
		return
	}

	m.Text = fmt.Sprintf("Task (%d) created", taskId)
	s.LastStep = s.LastStep.next
	s.task = nil
}

func (s *Service) CommandCompleteLast(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) {
	taskId, err := strconv.ParseInt(upd.Message.Text, 10, 64)
	m.Text = "Enter correct task ID"

	if err != nil {
		return
	}

	err = s.repo.DeleteTask(taskId)

	if err != nil {
		m.Text = fmt.Sprintf("There is no task with ID: %d", taskId)
	} else {
		m.Text = fmt.Sprintf("Task (%d) was completed!", taskId)
	}

	s.LastStep = s.LastStep.next

}
