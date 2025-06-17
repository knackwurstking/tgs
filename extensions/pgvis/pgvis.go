package pgvis

import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Data struct {
	Targets  *tgs.Targets          `yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `yaml:"register,omitempty"`
}

type PGVis struct {
	*tgbotapi.BotAPI

	data      *Data
	callbacks tgs.ReplyCallbacks
}

func New(api *tgbotapi.BotAPI) *PGVis {
	return &PGVis{
		data: &Data{
			Targets:  tgs.NewTargets(),
			Register: make([]tgs.BotCommandScope, 0),
		},
		callbacks: tgs.ReplyCallbacks{},
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (p *PGVis) Name() string {
	return "pgvis"
}

func (p *PGVis) SetBot(api *tgbotapi.BotAPI) {
	p.BotAPI = api
}

func (p *PGVis) ConfigPath() string {
	return "pgvis.yaml"
}

func (p *PGVis) MarshalYAML() (any, error) {
	return p.data, nil
}

func (p *PGVis) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(p.data)
}

func (p *PGVis) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/pgvisregister", "PG: Vis server registration", p.data.Register)
}

func (p *PGVis) Is(message *tgbotapi.Message) bool {
	return strings.HasPrefix(message.Command(), "pgvis")
}

func (p *PGVis) Handle(message *tgbotapi.Message) error {
	// TODO: ...

	return errors.New("under construction")
}
