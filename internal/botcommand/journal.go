package botcommand

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

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
		Register []tgs.BotCommandScope `json:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty"`
		Units    *Units                `json:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	})
}

func (this *Journal) MarshalYAML() (interface{}, error) {
	return struct {
		Register []tgs.BotCommandScope `yaml:"register,omitempty"`
		Targets  *Targets              `yaml:"targets,omitempty"`
		Units    *Units                `yaml:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	}, nil
}

func (this *Journal) UnmarshalJSON(data []byte) error {
	d := struct {
		Register []tgs.BotCommandScope `json:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty"`
		Units    *Units                `json:"units,omitempty"`
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
		Register []tgs.BotCommandScope `yaml:"register,omitempty"`
		Targets  *Targets              `yaml:"targets,omitempty"`
		Units    *Units                `yaml:"units,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
		Units:    this.units,
	}

	err := value.Decode(&d)
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
		return this.handleListCommand(message)
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
	return command == BotCommandJournal[1:]+"list"
}

func (this *Journal) handleListCommand(message *tgbotapi.Message) error {
	var (
		system []string
		user   []string
	)

	for _, unit := range this.units.System {
		system = append(system, fmt.Sprintf("<li>%s</li>", unit.Name))
	}

	for _, unit := range this.units.User {
		user = append(user, fmt.Sprintf("<li>%s</li>", unit.Name))
	}

	content := fmt.Sprintf(
		`<!doctype html><html>
			<head>
				<title>Journal Units</title>
			</head>

			<body>
				<h2>System Level</h2>
				<ul>
					%s
				</ul>

				<h2>User Level</h2>
				<ul>
					%s
				</ul>
			</body>
		</html>`,
		strings.Join(system, "\n"),
		strings.Join(user, "\n"),
	)

	documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  "journal-units.html",
		Bytes: []byte(content),
	})
	documentConfig.ReplyToMessageID = message.MessageID

	_, err := this.BotAPI.Send(documentConfig)
	return err
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

	var output string
	if unit.Output == "default" || unit.Output == "" {
		output = "short"
	} else {
		output = unit.Output
	}

	var cmd *exec.Cmd
	if isUser {
		cmd = exec.Command("journalctl",
			"--user",
			"-u", unit.Name,
			"--output", output,
			"--no-pager",
		)
	} else {
		cmd = exec.Command("journalctl",
			"-u", unit.Name,
			"--output", output,
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
