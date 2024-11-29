package config

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigStats struct {
	Register          []tgs.BotCommandScope `json:"register" yaml:"register"`
	ValidationTargets *ValidationTargets    `json:"targets" yaml:"targets"`
}

func NewCommandConfigStats() *CommandConfigStats {
	return &CommandConfigStats{
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
	}
}
