package botcommand

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type OPMangaConfig struct {
	Register []tgs.BotCommandScope `json:"register,omitempty"`
	Targets  *Targets              `json:"targets,omitempty"`
}

// OPManga implements the Handler interface
type OPManga struct {
	tgbotapi.BotAPI

	register []tgs.BotCommandScope
	targets  *Targets
}

func (this *OPManga) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *OPManga) Targets() *Targets {
	return this.targets
}

func (this *OPManga) Run(message *tgbotapi.Message) error {
	if this.isListCommand(message.Command()) {
		return this.handleListCommand(message)
	}

	return fmt.Errorf("under construction")
}

func (this *OPManga) AddCommands(c *tgs.MyBotCommands) {
	c.Add(BotCommandOPManga+"list", "List all available chapters", this.Register())
	c.Add(BotCommandOPManga, "Request a chapter", this.Register())
}

func (this *OPManga) MarshalJSON() ([]byte, error) {
	return json.Marshal(OPMangaConfig{Register: this.register, Targets: this.targets})
}

func (this *OPManga) UnmarshalJSON(data []byte) error {
	d := OPMangaConfig{Register: this.register, Targets: this.targets}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *OPManga) MarshalYAML() (interface{}, error) {
	return OPMangaConfig{Register: this.register, Targets: this.targets}, nil
}

func (this *OPManga) UnmarshalYAML(value *yaml.Node) error {
	d := OPMangaConfig{Register: this.register, Targets: this.targets}

	if err := value.Decode(&d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *OPManga) isListCommand(c string) bool {
	return c == BotCommandOPManga[1:]+"list"
}

func (this *OPManga) handleListCommand(m *tgbotapi.Message) error {
	buf := bytes.NewBuffer([]byte{})

	if t, err := template.ParseFS(Templates,
		"templates/index.html",
		"templates/opmangalist.html",
	); err != nil {
		return err
	} else {
		// TODO: Get a list with available chapters from the path, need a config field for this
		if err := t.Execute(buf, struct {
			PageTitle string
			// TODO: ...
		}{
			PageTitle: "One Piece Manga | Chapters",
			// TODO: ...
		}); err != nil {
			return err
		}
	}

	if content, err := io.ReadAll(buf); err != nil {
		return err
	} else {
		documentConfig := tgbotapi.NewDocument(m.Chat.ID, tgbotapi.FileBytes{
			Name:  "journal-units.html",
			Bytes: content,
		})
		documentConfig.ReplyToMessageID = m.MessageID

		_, err = this.BotAPI.Send(documentConfig)
		return err
	}
}
