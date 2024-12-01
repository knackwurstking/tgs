package config

import (
	"fmt"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigJournal struct {
	botcommand.Handler

	Register          []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
	ValidationTargets *ValidationTargets    `json:"targets,omitempty" yaml:"targets,omitempty"`
	Units             *Units                `json:"units,omitempty" yaml:"units,omitempty"`
}

func NewCommandConfigJournal(bot *tgbotapi.BotAPI) *CommandConfigJournal {
	return &CommandConfigJournal{
		Handler:           botcommand.NewJournal(bot),
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
		Units:             NewUnits(),
	}
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
