package botcommand

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/config"
)

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
		// TODO: Reply with a list with of available unit files

		return fmt.Errorf("under construction")
	}

	// TODO: Find out how to do this `tgbotapi.ReplyKeyboardMarkup` thing
	// 	- Number keyboard for entering the unit name

	return fmt.Errorf("under construction")
}

func (this *Journal) isListCommand(command string) bool {
	return strings.HasSuffix(command, config.BotCommandJournal+"list")
}
