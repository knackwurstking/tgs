package journal

import (
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type JournalConfig struct {
	Targets  *botcommand.Targets   `json:"targets,omitempty" yaml:"targets,omitempty"`
	Units    *Units                `json:"units,omitempty" yaml:"units,omitempty"`
	Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
}
