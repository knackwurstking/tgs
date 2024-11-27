package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Config struct {
	Token           string           `json:"token" yaml:"token"`
	IPCommandConfig *IPCommandConfig `json:"ip" yaml:"ip"`
}

func NewConfig() *Config {
	return &Config{
		IPCommandConfig: NewIPCommandConfig(),
	}
}

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

type ValidationsConfig struct {
	Users []UserValidation `json:"users,omitempty" yaml:"users,omitempty"`
	Chats []ChatValidation `json:"chats,omitempty" yaml:"chats,omitempty"`
}

func NewValidationsConfig() *ValidationsConfig {
	return &ValidationsConfig{
		Users: make([]UserValidation, 0),
		Chats: make([]ChatValidation, 0),
	}
}

type UserValidation struct {
	ID int64 `json:"id" yaml:"id"`
}

type ChatValidation struct {
	ID   int64  `json:"id" yaml:"id"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}
