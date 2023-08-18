package bot

import (
	"github.com/dkshi/hephaestus/internal/service"
	Botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	DebugMode     bool
	UpdateTimeout int
	UpdateOffset  int
}

type Bot struct {
	BotApi *Botapi.BotAPI
	srv    *service.Service
}

func NewBot(b *Botapi.BotAPI, s *service.Service) *Bot {
	return &Bot{
		BotApi: b,
		srv:    s,
	}
}

func (b *Bot) RunBot(config *Config) error {

	b.BotApi.Debug = config.DebugMode
	updateConfig := Botapi.NewUpdate(config.UpdateOffset)
	updateConfig.Timeout = config.UpdateTimeout

	updates := b.BotApi.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Handle next step if user previously used a command else wait for new command
		if b.srv.LastStep != nil {

			if err := b.srv.HandleNextStep(&update, &msg); err != nil {
				return err
			}
		} else {
			switch update.Message.Command() {
			case "create":
				b.srv.CommandCreate(&update, &msg)
			case "complete":
				msg.Text = "Choose a task to complete"
			case "profile":
				b.srv.CommandProfile(&update, &msg)
			default:
				msg.Text = "I don't know this command"
			}
		}

		if _, err := b.BotApi.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
