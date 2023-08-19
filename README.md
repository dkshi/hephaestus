# hephaestus: tg bot for task organization

<img align="left" width="300px" height="300px" src="https://i.ibb.co/S0Bp2wz/hammer.png">

This bot is written using <a href="https://github.com/go-telegram-bot-api/telegram-bot-api">tgbotapi</a> long polling and uses postgresDB to store data. Yes, I am the creator of a logo on the left side.

**The bot's functions include:** 

- Creating tasks
- Viewing previously created tasks
- Deleting (completing) tasks

## Getting started

Clone this repository and edit configs/config.yml and .env:
```yaml
db:
  username: "postgres"
  host: "localhost"
  port: "5436"
  dbname: "postgres"
  sslmode: "disable"
```

```.env
API_TOKEN=YOUR_TG_BOT_TOKEN
DB_PASSWORD=YOUR_DB_PASSWORD
```

And use go run command on main.go file:
```bash
$ go run cmd/bot/main.go
```

## Next step handler realisation

Since the next step handler is not implemented in the tgbotapi package, I have decided to use a linked list that will execute the current function in it and move to the next element until it becomes equal to nil.
```go
// Linked list for storing steps
type nextStepLinkedList struct {
	val  func(upd *tgbotapi.Update, m *tgbotapi.MessageConfig)
	next *nextStepLinkedList
}

// Function that will handle next step if LL not equal to nil
func (s *Service) HandleNextStep(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) error {
	if s.LastStep == nil {
		return errors.New("command linked list equal nil")
	}

	s.LastStep.val(upd, m)
	return nil
}

// The first command in a row needs to create a proper linked list
func (s *Service) CommandCreate(upd *tgbotapi.Update, m *tgbotapi.MessageConfig) {
	m.Text = "Enter a name for your new task"
	s.linkedCmdCreate(upd.Message.Chat.ID)
}

// Creating linked list with steps needed to procede (only for command /new)
func (s *Service) linkedCmdCreate(chatId int64) {
	s.LastStep = &nextStepLinkedList{
		val: s.CommandCreateSecond, // Some func that will handle next step
		next: &nextStepLinkedList{
			val:  s.CommandCreateLast, // Some func that will handle next step
			next: nil, // It was the last func
		},
	}

// Some other logic
	s.task = &hephaestus.Task{
		ChatID: chatId,
	}
}
```

## Long Polling
Here we obtain an updates channel and process each update in an infinite loop:
```go
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
				b.srv.CommandComplete(&update, &msg)
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
```



