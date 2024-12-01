package config

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func New(bot *tgbotapi.BotAPI) *Config {
	return NewConfig(bot)
}

type Config struct {
	Token   string                `json:"token" yaml:"token"`
	IP      *CommandConfigIP      `json:"ip,omitempty" yaml:"ip,omitempty"`
	Stats   *CommandConfigStats   `json:"stats,omitempty" yaml:"stats,omitempty"`
	Journal *CommandConfigJournal `json:"journal,omitempty" yaml:"journal,omitempty"`
}

func NewConfig(bot *tgbotapi.BotAPI) *Config {
	return &Config{
		IP:      NewCommandConfigIP(bot),
		Stats:   NewCommandConfigStats(bot),
		Journal: NewCommandConfigJournal(bot),
	}
}
