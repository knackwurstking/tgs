package config

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type IPCommandConfig struct {
	Register          []tgs.BotCommandScope `json:"register" yaml:"register"`
	ValidationsConfig *ValidationsConfig    `json:"targets" yaml:"targets"`
}

func NewIPCommandConfig() *IPCommandConfig {
	return &IPCommandConfig{
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationsConfig: NewValidationsConfig(),
	}
}
