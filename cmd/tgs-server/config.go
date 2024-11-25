package main

import (
	"fmt"

	"github.com/knackwurstking/tgs/pkg/data"
)

type Users []data.User
type Chats []data.Chat

type Config struct {
	Token    string         `json:"token"`
	Targets  Targets        `json:"targets"`
	Commands CommandConfigs `json:"commands"`
}

func NewConfig() *Config {
	return &Config{
		Targets: Targets{
			Users: make(Users, 0),
			Chats: make(Chats, 0),
		},
		Commands: CommandConfigs{},
	}
}

type Targets struct {
	Users Users `json:"users"`
	Chats Chats `json:"chats"`
}

type CommandConfigs struct {
	IP          CommandConfig `json:"ip"`
	JournalList CommandConfig `json:"journallist"`
	Journal     CommandConfig `json:"journal"`
	PicowStatus CommandConfig `json:"picowstatus"`
	PicowOn     CommandConfig `json:"picowon"`
	PicowOff    CommandConfig `json:"picowoff"`
	OPManga     CommandConfig `json:"opmanga"`
	OPMangaList CommandConfig `json:"opmangalist"`
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
	Targets *Targets `json:"targets"`
}
