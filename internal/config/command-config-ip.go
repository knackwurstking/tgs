package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type CommandConfigIP struct {
	botcommand.Handler

	Register          []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
	ValidationTargets *ValidationTargets    `json:"targets,omitempty" yaml:"targets,omitempty"`
}

// TODO: Add the `botcommand.IP` struct here
func NewCommandConfigIP(bot *tgbotapi.BotAPI) *CommandConfigIP {
	return &CommandConfigIP{
		Handler:           botcommand.NewIP(bot),
		Register:          make([]tgs.BotCommandScope, 0),
		ValidationTargets: NewValidationTargets(),
	}
}
