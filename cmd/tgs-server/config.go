package main

import (
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Config struct {
	Token    string      `json:"token" yaml:"token"`
	Targets  Targets     `json:"targets" yaml:"targets"`
	Commands BotCommands `json:"commands" yaml:"commands"`
}

type Targets struct {
	Users Users `json:"users" yaml:"users"`
	Chats Chats `json:"chats" yaml:"chats"`
}

type Users []*tgs.User

type Chats []*tgs.Chat

type BotCommands []BotCommand

type BotCommand struct {
	// TODO: ...
}
