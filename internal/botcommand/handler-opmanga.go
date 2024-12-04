package botcommand

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type OPMangaChapter struct {
	Name string
	Path string
}

func (this *OPMangaChapter) PDF() ([]byte, error) {
	// TODO: Read pdf data from path and return

	return nil, fmt.Errorf("under construction")
}

type OPMangaArc struct {
	Name     string
	Chapters []OPMangaChapter
}

type OPMangaTemplateData struct {
	PageTitle string
	Arcs      []OPMangaArc
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
	arcs, err := this.buildOPMangaArcs()
	if err != nil {
		return err
	}

	content, err := GetTemplateData(OPMangaTemplateData{
		PageTitle: "One Piece Manga | Chapters",
		Arcs:      arcs,
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

func (this *OPManga) buildOPMangaArcs() ([]OPMangaArc, error) {
	if this.path == "" {
		return nil, fmt.Errorf("missing path")
	}

	// TODO: Grep and build data from path

	return nil, fmt.Errorf("under construction")
}
