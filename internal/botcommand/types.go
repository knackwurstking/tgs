package botcommand

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Handler interface {
	Run(message *tgbotapi.Message) error
	AddCommands(c *tgs.MyBotCommands, scopes ...tgs.BotCommandScope)
}
