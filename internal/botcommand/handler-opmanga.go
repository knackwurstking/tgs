package botcommand

import (
	"encoding/json"

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

func (this *OPManga) Run(message *tgbotapi.Message) error

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
