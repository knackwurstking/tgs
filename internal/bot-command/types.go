package botcommand

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler interface {
	Run(message *tgbotapi.Message) error
}
