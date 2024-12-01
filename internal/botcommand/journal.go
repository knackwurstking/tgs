package botcommand

import (
	"encoding/json"
	"fmt"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Journal struct {
	*tgbotapi.BotAPI

	register []tgs.BotCommandScope
	targets  *Targets
	units    *Units
}

func NewJournal(botAPI *tgbotapi.BotAPI) *Journal {
	return &Journal{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  NewTargets(),
		units:    NewUnits(),
	}
}

func (this *Journal) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
		Units    *Units                `json:"units,omitempty" yaml:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	})
}

func (this *Journal) MarshalYAML() (interface{}, error) {
	return struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
		Units    *Units                `json:"units,omitempty" yaml:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	}, nil
}

func (this *Journal) UnmarshalJSON(data []byte) error {
	d := struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
		Units    *Units                `json:"units,omitempty" yaml:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	}

	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets
	this.units = d.Units

	return nil
}

func (this *Journal) UnmarshalYAML(value *yaml.Node) error {
	d := struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
		Units    *Units                `json:"units,omitempty" yaml:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	}

	err := value.Encode(&d)
	if err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets
	this.units = d.Units

	return nil
}

func (this *Journal) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *Journal) Targets() *Targets {
	return this.targets
}

func (this *Journal) Run(message *tgbotapi.Message) error {
	if this.isListCommand(message.Command()) {
		// TODO: Maybe create a html file here? using some templating or whatever?
		content := "System Level\n\n"

		for _, unit := range this.units.System {
			content += fmt.Sprintf("\t- %s\n", unit.Name)
		}

		content += "\nUser Level\n\n"

		for _, unit := range this.units.User {
			content += fmt.Sprintf("\t- %s\n", unit.Name)
		}

		documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
			Name:  "journal-units-list.txt",
			Bytes: []byte(content),
		})
		documentConfig.ReplyToMessageID = message.MessageID

		_, err := this.BotAPI.Send(documentConfig)
		return err
	}

	// TODO: Find out how to do this `tgbotapi.ReplyKeyboardMarkup` thing
	// 	- Number keyboard for entering the unit name

	return fmt.Errorf("under construction")
}

func (this *Journal) AddCommands(c *tgs.MyBotCommands, scopes ...tgs.BotCommandScope) {
	c.Add(BotCommandJournal+"list", "List journalctl logs", scopes)
	c.Add(BotCommandJournal, "Get a journalctl log", scopes)
}

func (this *Journal) isListCommand(command string) bool {
	return command == BotCommandJournal+"list"
}

type Units struct {
	System []Unit `json:"system,omitempty" yaml:"system,omitempty"`
	User   []Unit `json:"user,omitempty" yaml:"user,omitempty"`
}

func NewUnits() *Units {
	return &Units{
		System: make([]Unit, 0),
		User:   make([]Unit, 0),
	}
}

func (this *Units) GetSystemUnit(name string) (*Unit, error) {
	for i, u := range this.System {
		if u.Name == name {
			return &this.System[i], nil
		}
	}

	return nil, fmt.Errorf("user unit %s not found", name)
}

func (this *Units) GetUserUnit(name string) (*Unit, error) {
	for i, u := range this.User {
		if u.Name == name {
			return &this.User[i], nil
		}
	}

	return nil, fmt.Errorf("system unit %s not found", name)
}

func (this *Units) GetOutputForUnitPerName(name string) (data []byte, err error) {
	isUser := true

	unit, err := this.GetUserUnit(name)
	if err != nil {
		isUser = false

		unit, err = this.GetSystemUnit(name)
		if err != nil {
			return nil, err
		}
	}

	var cmd *exec.Cmd
	if isUser {
		cmd = exec.Command("journalctl",
			"--user",
			"-u", unit.Name,
			"--output", unit.GetOutput(),
			"--no-pager",
		)
	} else {
		cmd = exec.Command("journalctl",
			"-u", unit.Name,
			"--output", unit.GetOutput(),
			"--no-pager",
		)
	}

	if data, err = cmd.CombinedOutput(); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

type Unit struct {
	// Name or the unit file to get
	Name string `json:"name" yaml:"name"`
	// Output will be used for the shell command `journalctl` as `--output ${output}`
	//
	// optional
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
}

// GetOutput will just return the `this.Output` field, but will do some parsing if the value
// is empty or "default"
func (this *Unit) GetOutput() string {
	output := this.Output
	if output == "" || output == "default" {
		return "short"
	}

	return this.Output
}
