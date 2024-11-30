package config

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigStats struct {
	Register          []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
	ValidationTargets *ValidationTargets    `json:"targets,omitempty" yaml:"targets,omitempty"`
}

func NewCommandConfigStats() *CommandConfigStats {
	return &CommandConfigStats{
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
	}
}
