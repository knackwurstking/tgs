package config

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigIP struct {
	Register          []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
	ValidationTargets *ValidationTargets    `json:"targets,omitempty" yaml:"targets,omitempty"`
}

func NewCommandConfigIP() *CommandConfigIP {
	return &CommandConfigIP{
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
	}
}
