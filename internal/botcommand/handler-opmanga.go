package botcommand

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type OPMangaTemplateData struct {
	PageTitle string
	Chapters  OPMangaChapters
}

type OPMangaChapters struct {
	// ...
}

func NewOPMangaChapters() *OPMangaChapters {
	return &OPMangaChapters{}
}

func (this *OPMangaChapters) Grep(path string) error {
	// TODO: ...

	return fmt.Errorf("under construction")
}

type OPMangaConfig struct {
	Register []tgs.BotCommandScope `json:"register,omitempty"`
	Targets  *Targets              `json:"targets,omitempty"`
	Path     string                `json:"path" json:"path"`
}

// OPManga implements the Handler interface
type OPManga struct {
	tgbotapi.BotAPI

	register []tgs.BotCommandScope
	targets  *Targets
	path     string
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

	// TODO: ...

	return fmt.Errorf("under construction")
}

func (this *OPManga) AddCommands(c *tgs.MyBotCommands) {
	c.Add(BotCommandOPManga+"list", "List all available chapters", this.Register())
	c.Add(BotCommandOPManga, "Request a chapter", this.Register())
}

func (this *OPManga) MarshalJSON() ([]byte, error) {
	return json.Marshal(OPMangaConfig{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	})
}

func (this *OPManga) UnmarshalJSON(data []byte) error {
	d := OPMangaConfig{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *OPManga) MarshalYAML() (interface{}, error) {
	return OPMangaConfig{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}, nil
}

func (this *OPManga) UnmarshalYAML(value *yaml.Node) error {
	d := OPMangaConfig{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}

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
	content, err := GetTemplateData(OPMangaTemplateData{
		PageTitle: "One Piece Manga | Chapters",
	})
	if err != nil {
		return err
	}

	documentConfig := tgbotapi.NewDocument(m.Chat.ID, tgbotapi.FileBytes{
		Name:  "journal-units.html",
		Bytes: content,
	})
	documentConfig.ReplyToMessageID = m.MessageID

	_, err = this.BotAPI.Send(documentConfig)
	return err
}

func (this *OPManga) getChapters() (*OPMangaChapters, error) {
	chapters := NewOPMangaChapters()

	if err := chapters.Grep(this.path); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("under construction")
}
