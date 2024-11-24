package main

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Users []*tgs.User
type Chats []*tgs.Chat

type Config struct {
	Token    string         `json:"token" yaml:"token"`
	Targets  Targets        `json:"targets" yaml:"targets"`
	Commands CommandConfigs `json:"commands" yaml:"commands"`
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
	Users Users `json:"users" yaml:"users"`
	Chats Chats `json:"chats" yaml:"chats"`
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

type CommandConfig struct {
	Targets *Targets `json:"targets" yaml:"targets"`
}
