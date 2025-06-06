package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/botcommand/ip"
	"github.com/knackwurstking/tgs/internal/botcommand/journal"
	"github.com/knackwurstking/tgs/internal/botcommand/opmanga"
	"github.com/knackwurstking/tgs/internal/botcommand/stats"
)

func New(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Config {
	return NewConfig(bot, reply)
}

type Config struct {
	IP      *ip.IP                 `json:"ip,omitempty" yaml:"ip,omitempty"`
	Stats   *stats.Stats           `json:"stats,omitempty" yaml:"stats,omitempty"`
	Journal *journal.Journal       `json:"journal,omitempty" yaml:"journal,omitempty"`
	OPManga *opmanga.OPManga       `json:"opmanga,omitempty" yaml:"opmanga,omitempty" `
	Reply   chan *botcommand.Reply `json:"-" yaml:"-"`
	Token   string                 `json:"token" yaml:"token"`
}

func NewConfig(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Config {
	return &Config{
		IP:      ip.NewIP(bot),
		Stats:   stats.NewStats(bot),
		Journal: journal.NewJournal(bot, reply),
		OPManga: opmanga.NewOPManga(bot, reply),
		Reply:   reply,
	}
}

func (this *Config) SetBot(bot *tgbotapi.BotAPI) {
	this.IP.BotAPI = bot
	this.Stats.BotAPI = bot
	this.Journal.BotAPI = bot
	this.OPManga.BotAPI = bot
}
