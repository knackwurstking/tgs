package config

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigIP struct {
	Register          []tgs.BotCommandScope `json:"register" yaml:"register"`
	ValidationsConfig *ValidationTargets    `json:"targets" yaml:"targets"`
}

func NewCommandConfigIP() *CommandConfigIP {
	return &CommandConfigIP{
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationsConfig: NewValidationsConfig(),
	}
}
