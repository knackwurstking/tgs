package config

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type IPCommandConfig struct {
	Register          []tgbotapi.BotCommandScope `json:"register" yaml:"register"`
	ValidationsConfig *ValidationsConfig         `json:"targets" yaml:"targets"`
}

func NewIPCommandConfig() *IPCommandConfig {
	return &IPCommandConfig{
		Register:          make([]tgbotapi.BotCommandScope, 0),
		ValidationsConfig: NewValidationsConfig(),
	}
}
