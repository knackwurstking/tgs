package botcommand

import (
	"embed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

const (
	BotCommandIP      string = "/ip"
	BotCommandStats   string = "/stats"
	BotCommandJournal string = "/journal"
	BotCommandOPManga string = "/opmanga"
)

var (
	//go:embed templates
	Templates embed.FS
)

type Handler interface {
	Register() []tgs.BotCommandScope
	Targets() *Targets
	Run(message *tgbotapi.Message) error
	AddCommands(c *tgs.MyBotCommands)

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error

	MarshalYAML() (interface{}, error)
	UnmarshalYAML(value *yaml.Node) error
}
