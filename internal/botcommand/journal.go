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
	// Check command for a sub command like "list" first
	if strings.HasSuffix(message.Command(), config.BotCommandJournal+"list") {
		// ...

		return fmt.Errorf("under construction")
	}

	// TODO: Find out how to do this `tgbotapi.ReplyKeyboardMarkup` thing
	// 	- Number keyboard to entering the episode number

	return fmt.Errorf("under construction")
}
