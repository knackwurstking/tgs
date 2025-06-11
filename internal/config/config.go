package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/botcommand/opmanga"
	"github.com/knackwurstking/tgs/internal/botcommand/stats"
)

func New(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Config {
	return NewConfig(bot, reply)
}

type Config struct {
	Stats   *stats.Stats           `json:"stats,omitempty" yaml:"stats,omitempty"`
	OPManga *opmanga.OPManga       `json:"opmanga,omitempty" yaml:"opmanga,omitempty" `
	Reply   chan *botcommand.Reply `json:"-" yaml:"-"`
	Token   string                 `json:"token" yaml:"token"`
}

func NewConfig(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Config {
	return &Config{
		Stats:   stats.NewStats(bot),
		OPManga: opmanga.NewOPManga(bot, reply),
		Reply:   reply,
	}
}

func (this *Config) SetBot(bot *tgbotapi.BotAPI) {
	this.Stats.BotAPI = bot
	this.OPManga.BotAPI = bot
}
