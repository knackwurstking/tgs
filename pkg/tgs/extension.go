package tgs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

type Extension interface {
	Name() string
	SetBot(api *tgbotapi.BotAPI)
	ConfigPath() string
	MarshalYAML() (any, error)
	UnmarshalYAML(value *yaml.Node) error
	Commands(mbc *MyBotCommands)
	Is(message *tgbotapi.Message) bool
	Handle(message *tgbotapi.Message) error
}
