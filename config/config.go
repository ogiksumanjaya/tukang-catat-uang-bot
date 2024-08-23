package config

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

type Config struct{}

func (c *Config) GetDatabaseConfig() string {
	return os.Getenv("TGBOT_DB_POSTGRESQL_URL")
}

func (c *Config) GetServerAddress() string {
	restPort, err := strconv.Atoi(os.Getenv("SERVER_REST_PORT"))
	if err != nil {
		log.Println("SERVER_REST_PORT not found, using default port 8080")
	}
	return fmt.Sprintf("%s:%d", os.Getenv("SERVER_HOST"), restPort)
}

type Telegram struct {
	Token           string
	UsernameAllowed string
}

func (c *Config) GetTelegram() Telegram {
	return Telegram{
		Token:           os.Getenv("TELEGRAM_BOT_TOKEN"),
		UsernameAllowed: os.Getenv("TELEGRAM_USERNAME_ALLOWED"),
	}
}

func (c *Config) InitConfigTelegram() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(c.GetTelegram().Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	return bot, updates, nil
}
