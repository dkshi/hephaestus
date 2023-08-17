package service

import (
	"errors"

	"github.com/dkshi/hephaestus"
	"github.com/dkshi/hephaestus/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	LastStep *nextStepLinkedList
	repo *repository.Repository
	task *hephaestus.Task
}

// Linked list for storing steps
type nextStepLinkedList struct {
	val  func(upd *tgbotapi.Update, m *tgbotapi.MessageConfig)
	next *nextStepLinkedList
}

func NewService(r *	repository.Repository) *Service {
	return &Service{
		LastStep: nil,
		repo: r,
		task: nil,
	}
}

func (s *Service) HandleNextStep(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) error {
	if s.LastStep == nil {
		return errors.New("command linked list equal nil")
	}

	s.LastStep.val(upd, m)
	return nil
}

func (s *Service) CommandCreate(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) {
	m.Text = "Enter a name for your new task"
	s.linkedCmdCreate(upd.Message.Chat.ID)
}
