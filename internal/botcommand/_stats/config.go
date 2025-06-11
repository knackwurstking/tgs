package stats

import (
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type StatsConfig struct {
	Targets  *botcommand.Targets   `json:"targets,omitempty" yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
}
