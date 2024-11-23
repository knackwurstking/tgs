package main

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Users []*tgs.User
type Chats []*tgs.Chat

type Config struct {
	Token    string      `json:"token" yaml:"token"`
	Targets  Targets     `json:"targets" yaml:"targets"`
	Commands BotCommands `json:"commands" yaml:"commands"`
}

func NewConfig() *Config {
	return &Config{
		Targets: Targets{
			Users: make(Users, 0),
			Chats: make(Chats, 0),
		},
		Commands: BotCommands{},
	}
}

type Targets struct {
	Users Users `json:"users" yaml:"users"`
	Chats Chats `json:"chats" yaml:"chats"`
}

type BotCommands struct {
	IP          BotCommand `json:"ip" yaml:"ip"`
	JournalList BotCommand `json:"journallist" yaml:"journallist"`
	Journal     BotCommand `json:"journal" yaml:"journal"`
	PicowStatus BotCommand `json:"picowstatus" yaml:"picowstatus"`
	PicowOn     BotCommand `json:"picowon" yaml:"picowon"`
	PicowOff    BotCommand `json:"picowoff" yaml:"picowoff"`
	OPManga     BotCommand `json:"opmanga" yaml:"opmanga"`
	OPMangaList BotCommand `json:"opmangalist" yaml:"opmangalist"`
}

type BotCommand struct {
	Targets *Targets `json:"targets" yaml:"targets"`
}
