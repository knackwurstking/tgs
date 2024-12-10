package opmanga

import (
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Config struct {
	Targets  *botcommand.Targets   `json:"targets,omitempty" yaml:"targets,omitempty"`
	Path     string                `json:"path" yaml:"path"`
	Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
}
