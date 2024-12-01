package botcommand

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Handler interface {
	MarshalJSON() ([]byte, error)
	MarshalYAML() (interface{}, error)
	UnmarshalJSON(data []byte) error
	UnmarshalYAML(value *yaml.Node) error
	Register() []tgs.BotCommandScope
	Targets() *ValidationTargets
	Run(message *tgbotapi.Message) error
	AddCommands(c *tgs.MyBotCommands, scopes ...tgs.BotCommandScope)
}
