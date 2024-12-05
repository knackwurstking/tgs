package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
)

func New(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Config {
	return NewConfig(bot, reply)
}

type Config struct {
	Token   string              `json:"token" yaml:"token"`
	IP      *botcommand.IP      `json:"ip,omitempty" yaml:"ip,omitempty"`
	Stats   *botcommand.Stats   `json:"stats,omitempty" yaml:"stats,omitempty"`
	Journal *botcommand.Journal `json:"journal,omitempty" yaml:"journal,omitempty"`
	OPManga *botcommand.OPManga `json:"opmanga,omitempty" yaml:"opmanga,omitempty" `

	Reply chan *botcommand.Reply `json:"-" yaml:"-"`
}

func NewConfig(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Config {
	return &Config{
		IP:      botcommand.NewIP(bot),
		Stats:   botcommand.NewStats(bot),
		Journal: botcommand.NewJournal(bot, reply),
		OPManga: botcommand.NewOPManga(bot),
		Reply:   reply,
	}
}
