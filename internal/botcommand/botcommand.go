package botcommand

import (
	"bytes"
	"embed"
	"html/template"
	"io"

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

//go:embed templates
var Templates embed.FS

type Handler interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error

	MarshalYAML() (interface{}, error)
	UnmarshalYAML(value *yaml.Node) error

	Register() []tgs.BotCommandScope
	Targets() *Targets
	AddCommands(c *tgs.MyBotCommands)
	Run(message *tgbotapi.Message) error
}

type TemplateData interface {
	Patterns() []string
}

type File interface {
	Name() string
	Data() []byte
}

func GetTemplateData(templateData TemplateData) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if t, err := template.ParseFS(Templates, templateData.Patterns()...); err != nil {
		return nil, err
	} else {
		if err := t.Execute(buf, templateData); err != nil {
			return nil, err
		}
	}

	return io.ReadAll(buf)
}
