package opmanga

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

// OPManga implements the Handler interface
type OPManga struct {
	*tgbotapi.BotAPI
	targets  *botcommand.Targets
	reply    chan *botcommand.Reply
	path     string
	register []tgs.BotCommandScope
}

func NewOPManga(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *OPManga {
	return &OPManga{
		BotAPI: bot,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),

		reply: reply,
	}
}

func (opm *OPManga) MarshalJSON() ([]byte, error) {
	return json.Marshal(Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	})
}

func (opm *OPManga) UnmarshalJSON(data []byte) error {
	d := Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	opm.register = d.Register
	opm.targets = d.Targets
	opm.path = d.Path

	return nil
}

func (opm *OPManga) MarshalYAML() (interface{}, error) {
	return Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	}, nil
}

func (opm *OPManga) UnmarshalYAML(value *yaml.Node) error {
	d := Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	}

	if err := value.Decode(&d); err != nil {
		return err
	}

	opm.register = d.Register
	opm.targets = d.Targets
	opm.path = d.Path

	return nil
}

func (opm *OPManga) BotCommand() string {
	return "opmanga"
}

func (opm *OPManga) Register() []tgs.BotCommandScope {
	return opm.register
}

func (opm *OPManga) Targets() *botcommand.Targets {
	return opm.targets
}

func (opm *OPManga) AddCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/"+opm.BotCommand()+"list", "List all available chapters", opm.Register())
	mbc.Add("/"+opm.BotCommand(), "Request a chapter", opm.Register())
}

func (opm *OPManga) Run(m *tgbotapi.Message) error {
	if opm.isListCommand(m.Command()) {
		return opm.handleListCommand(m)
	}
}

func (opm *OPManga) isListCommand(c string) bool {
	return c == opm.BotCommand()+"list"
}
