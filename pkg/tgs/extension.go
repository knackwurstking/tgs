package tgs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

type Extension interface {
	Name() string                         // Name is used for logging
	SetBot(api *tgbotapi.BotAPI)          // SetBot will pass the bot api to the extension
	ConfigPath() string                   // ConfigPath gets the path to the extension config, will be joined with the base config path
	MarshalYAML() (any, error)            // MarshalYAML config from extension struct
	UnmarshalYAML(value *yaml.Node) error // UnmarshalYAML config to extension struct
	AddBotCommands(mbc *MyBotCommands)    // AddBotCommands `mbc.Add(...)`
	Is(update tgbotapi.Update) bool       // Is checks if message belongs to this extension
	Handle(update tgbotapi.Update) error  // Handle will do shit
}
