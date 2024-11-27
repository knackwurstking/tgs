package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCommands struct {
	Commands map[string][]tgbotapi.BotCommand
}

func NewBotCommands() *BotCommands {
	return &BotCommands{
		Commands: map[string][]tgbotapi.BotCommand{},
	}
}

func (this *BotCommands) Add(command string, description string, scopes []tgbotapi.BotCommandScope) {
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

func (this *BotCommands) Register(bot *tgbotapi.BotAPI) error {
	for scope, botCommands := range this.Commands {
		scopeSplit := strings.SplitN(scope, ":", 3)

		scopeType := scopeSplit[0]
		scopeChatID, _ := strconv.ParseInt(scopeSplit[1], 10, 64)
		scopeUserID, _ := strconv.ParseInt(scopeSplit[2], 10, 64)

		setMyCommandsConfig := tgbotapi.NewSetMyCommands(botCommands...)
		setMyCommandsConfig.Scope = &tgbotapi.BotCommandScope{
			Type:   scopeType,
			ChatID: int64(scopeChatID),
			UserID: int64(scopeUserID),
		}

		_, err := bot.Send(setMyCommandsConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
