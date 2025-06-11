package journal

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

// Journal implements the Handler interface
type Journal struct {
	*tgbotapi.BotAPI
	targets  *botcommand.Targets
	units    *Units
	reply    chan *botcommand.Reply
	register []tgs.BotCommandScope
}

func NewJournal(botAPI *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Journal {
	return &Journal{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),
		units:    NewUnits(),

		reply: reply,
	}
}

func (j *Journal) MarshalJSON() ([]byte, error) {
	return json.Marshal(JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	})
}

func (j *Journal) UnmarshalJSON(data []byte) error {
	d := JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	j.register = d.Register
	j.targets = d.Targets
	j.units = d.Units

	return nil
}

func (j *Journal) MarshalYAML() (interface{}, error) {
	return JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	}, nil
}

func (j *Journal) UnmarshalYAML(value *yaml.Node) error {
	d := JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	}

	if err := value.Decode(&d); err != nil {
		return err
	}

	j.register = d.Register
	j.targets = d.Targets
	j.units = d.Units

	return nil
}
