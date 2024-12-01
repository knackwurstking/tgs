package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigStats struct {
	botcommand.Handler

	Register          []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
	ValidationTargets *ValidationTargets    `json:"targets,omitempty" yaml:"targets,omitempty"`
}

func NewCommandConfigStats(bot *tgbotapi.BotAPI) *CommandConfigStats {
	return &CommandConfigStats{
		Handler:           botcommand.NewStats(bot),
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
	}
}
