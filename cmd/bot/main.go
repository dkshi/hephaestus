package main

import (
	"os"

	"github.com/dkshi/hephaestus/internal/bot"
	"github.com/dkshi/hephaestus/internal/repository"
	"github.com/dkshi/hephaestus/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// Initializing config.yml
	if err := initConfig(); err != nil {
		logrus.Fatalf("error loading configs: %s", err.Error())
	}

	// Loading .env
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading .env: %s", err.Error())
	}

	//Initializing Bot API from package tgbotapi
	botapi, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		logrus.Fatalf("error creating bot api: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("error db connecting: %s", err.Error())
	}

	r := repository.NewRepository(db)
	r.CreateTables()
	h := service.NewService(r)
	b := bot.NewBot(botapi, h)

	b.RunBot(&bot.Config{
		DebugMode:     true,
		UpdateOffset:  0,
		UpdateTimeout: 30,
	})
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
