package config

import (
	"fmt"

	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigJournal struct {
	Register          []tgs.BotCommandScope `json:"register" yaml:"register"`
	ValidationTargets *ValidationTargets    `json:"targets" yaml:"targets"`
}

func NewCommandConfigJournal() *CommandConfigJournal {
	return &CommandConfigJournal{
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
	}
}

type Units struct {
	System []string `json:"system" yaml:"system"`
	User   []string `json:"user" yaml:"user"`
}

func (this *Units) Get(unit string, fromUser bool) (data []byte, err error) {
	data = make([]byte, 0)

	units := make([]string, 0)
	if fromUser {
		units = append(units, this.User...)
	} else {
		units = append(units, this.System...)
	}

	// TODO: Get journal from unit from user or system

	return nil, fmt.Errorf("under construction")
}
