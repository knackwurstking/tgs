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

func GetTemplateData(templateData any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if t, err := template.ParseFS(Templates,
		"templates/index.html",
		"templates/opmangalist.html",
	); err != nil {
		return nil, err
	} else {
		if err := t.Execute(buf, templateData); err != nil {
			return nil, err
		}
	}

	return io.ReadAll(buf)
}
