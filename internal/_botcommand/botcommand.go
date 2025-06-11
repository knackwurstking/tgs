package botcommand

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Handler interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error

	MarshalYAML() (interface{}, error)
	UnmarshalYAML(value *yaml.Node) error

	BotCommand() string

	Register() []tgs.BotCommandScope
	Targets() *Targets
	AddCommands(mbc *tgs.MyBotCommands)
	Run(message *tgbotapi.Message) error
}

type File interface {
	Name() string
	Data() []byte
}
