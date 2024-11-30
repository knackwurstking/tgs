package config

import (
	"fmt"

	"github.com/knackwurstking/tgs/pkg/tgs"
)

// TODO: Add Units ("units") field
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

func NewUnits() *Units {
	return &Units{
		System: make([]string, 0),
		User:   make([]string, 0),
	}
}

func (this *Units) Has(unit string) bool {
	units := make([]string, 0)
	units = append(units, this.System...)
	units = append(units, this.User...)
	for _, u := range this.System {
		if u == unit {
			return true
		}
	}

	return false
}

func (this *Units) IsSystemUnit(unit string) bool {
	for _, u := range this.System {
		if u == unit {
			return true
		}
	}

	return false
}

func (this *Units) IsUserUnit(unit string) bool {
	for _, u := range this.User {
		if u == unit {
			return true
		}
	}

	return false
}

func (this *Units) Get(unit string) (data []byte, err error) {
	var isUser bool
	if this.IsUserUnit(unit) {
		isUser = true
	} else if !this.IsSystemUnit(unit) {
		return nil, fmt.Errorf("missing unit: %s", unit)
	}

	data = make([]byte, 0)

	if isUser {
		// TODO: Create user command here for journal
	} else {
		// TODO: Create system command here for journal
	}

	// TODO: Run command here

	return nil, fmt.Errorf("under construction")
}
