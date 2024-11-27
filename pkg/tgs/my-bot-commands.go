package tgs

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MyBotCommands struct {
	Commands map[string][]tgbotapi.BotCommand
}

func NewMyBotCommands() *MyBotCommands {
	return &MyBotCommands{
		Commands: map[string][]tgbotapi.BotCommand{},
	}
}

func (this *MyBotCommands) Add(command string, description string, scopes []BotCommandScope) {
	for _, scope := range scopes {
		scopeString := fmt.Sprintf("%s:%d:%d", scope.Type, scope.ChatID, scope.UserID)

		if _, ok := this.Commands[scopeString]; !ok {
			this.Commands[scopeString] = []tgbotapi.BotCommand{}
		}

		this.Commands[scopeString] = append(this.Commands[scopeString], tgbotapi.BotCommand{
			Command:     command,
			Description: description,
		})
	}
}

func (this *MyBotCommands) Register(bot *tgbotapi.BotAPI) error {
	for scope, botCommands := range this.Commands {
		scopeSplit := strings.SplitN(scope, ":", 3)

		scopeType := scopeSplit[0]
		scopeChatID, _ := strconv.ParseInt(scopeSplit[1], 10, 64)
		scopeUserID, _ := strconv.ParseInt(scopeSplit[2], 10, 64)

		setMyCommandsConfig := tgbotapi.NewSetMyCommands(botCommands...)
		setMyCommandsConfig.Scope = &tgbotapi.BotCommandScope{
			Type:   scopeType,
			ChatID: scopeChatID,
			UserID: scopeUserID,
		}

		_, err := bot.Request(setMyCommandsConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
