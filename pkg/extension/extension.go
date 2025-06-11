package extension

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Extension interface {
	SetBot(api *tgbotapi.BotAPI)
	ConfigPath() string
	MarshalYAML() (any, error)
	UnmarshalYAML(value *yaml.Node) error
	Register() []tgs.BotCommandScope
	Targets() *botcommand.Targets
	Commands(mbc *tgs.MyBotCommands)
	Is(message *tgbotapi.Message) bool
	Handle(message *tgbotapi.Message) error
}
