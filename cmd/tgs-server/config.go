package main

import (
	"fmt"

	"github.com/knackwurstking/tgs/pkg/data"
)

type Users []data.User
type Chats []data.Chat

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

func (c *CommandConfigs) Get(command string) (*CommandConfig, error) {
	switch command {
	case "/ip":
		return &c.IP, nil
	case "/journallist":
		return &c.JournalList, nil
	case "/journal":
		return &c.Journal, nil
	case "/picowstatus":
		return &c.PicowStatus, nil
	case "/picowon":
		return &c.PicowOn, nil
	case "/picowoff":
		return &c.PicowOff, nil
	case "/opmanga":
		return &c.OPManga, nil
	case "/opmangalist":
		return &c.OPMangaList, nil
	default:
		return nil, fmt.Errorf("command %s not found", command)
	}
}

type CommandConfig struct {
	Targets *Targets `json:"targets" yaml:"targets"`
}
