package extension

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Extension interface {
	ConfigPath() string
	MarshalYAML() (any, error)
	UnmarshalYAML(value *yaml.Node) error
	Register() []tgs.BotCommandScope
	Targets() *botcommand.Targets
	Commands(mbc *tgs.MyBotCommands)
	Handler(message *tgbotapi.Message) error
}
