package botcommand

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

// TODO: Combine this with the `config.CommandConfigJournal` struct
type Journal struct {
	*tgbotapi.BotAPI
}

func NewJournal(botAPI *tgbotapi.BotAPI) *Journal {
	return &Journal{
		BotAPI: botAPI,
	}
}

func (this *Journal) Run(message *tgbotapi.Message) error {
	if this.isListCommand(message.Command()) {
		// TODO: Reply with a list of available unit files

		return fmt.Errorf("under construction")
	}

	// TODO: Find out how to do this `tgbotapi.ReplyKeyboardMarkup` thing
	// 	- Number keyboard for entering the unit name

	return fmt.Errorf("under construction")
}

func (this *Journal) AddCommands(c *tgs.MyBotCommands, scopes ...tgs.BotCommandScope) {
	c.Add(BotCommandJournal+"list", "List journalctl logs", scopes)
	c.Add(BotCommandJournal, "Get a journalctl log", scopes)
}

func (this *Journal) isListCommand(command string) bool {
	return command == BotCommandJournal+"list"
}
