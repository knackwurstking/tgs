package main

import (
	"fmt"
)

type Users []struct {
	ID int `json:"id" yaml:"id"`
}

type Chats []struct {
	ID   int    `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

type Config struct {
	Token    string         `json:"token" yaml:"token"`
	Commands CommandConfigs `json:"commands" yaml:"commands"`
}

func NewConfig() *Config {
	return &Config{
		Commands: CommandConfigs{},
	}
}

type CommandConfigs struct {
	IP          CommandConfig `json:"ip" yaml:"ip"`
	JournalList CommandConfig `json:"journallist" yaml:"journallist"`
	Journal     CommandConfig `json:"journal" yaml:"journal"`
	PicowStatus CommandConfig `json:"picowstatus" yaml:"picowstatus"`
	PicowOn     CommandConfig `json:"picowon" yaml:"picowon"`
	PicowOff    CommandConfig `json:"picowoff" yaml:"picowoff"`
	OPManga     CommandConfig `json:"opmanga" yaml:"opmanga"`
	OPMangaList CommandConfig `json:"opmangalist" yaml:"opmangalist"`
}

func (c *CommandConfigs) Get(command string) (*CommandConfig, error) {
	switch command {
	case BotCommandIP:
		return &c.IP, nil

	case BotCommandJournalList:
		return &c.JournalList, nil
	case BotCommandJournal:
		return &c.Journal, nil

	case BotCommandPicowStatus:
		return &c.PicowStatus, nil
	case BotCommandPicowON:
		return &c.PicowOn, nil
	case BotCommandPicowOFF:
		return &c.PicowOff, nil

	case BotCommandOPManga:
		return &c.OPManga, nil
	case BotCommandOPMangaList:
		return &c.OPMangaList, nil
	default:
		return nil, fmt.Errorf("command %s not found", command)
	}
}

type CommandConfig struct {
	Targets  *Targets `json:"targets,omitempty" yaml:"targets,omitempty"`   // [Optional]
	Disabled bool     `json:"disabled,omitempty" json:"disabled,omitempty"` // [Optional]
	Scopes   []Scope  `json:"scopes" yaml:"scopes"`                         // [Optional]
}

type Scope struct {
	Scope  string `json:"scope" yaml:"scope"`
	ChatID int    `json:"chat_id" yaml:"chat_id"`
	UserID int    `json:"user_id" yaml:"user_id"`
}

type Targets struct {
	Users Users `json:"users" yaml:"users"`
	Chats Chats `json:"chats" yaml:"chats"`
}
