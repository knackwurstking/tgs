package tgs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
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

func (this *MyBotCommands) Add(command string, description string, scopes []Scope) {
	log.Debug("Adding bot command",
		"command", command,
		"description", description,
		"scopes_count", len(scopes),
	)

	for i, scope := range scopes {
		scopeString := fmt.Sprintf("%s:%d:%d", scope.Type, scope.ChatID, scope.UserID)

		log.Debug("Processing command scope",
			"command", command,
			"scope_index", i+1,
			"scope_type", scope.Type,
			"chat_id", scope.ChatID,
			"user_id", scope.UserID,
		)

		if _, ok := this.Commands[scopeString]; !ok {
			this.Commands[scopeString] = []tgbotapi.BotCommand{}
			log.Debug("Created new scope command list", "scope", scopeString)
		}

		this.Commands[scopeString] = append(this.Commands[scopeString], tgbotapi.BotCommand{
			Command:     command,
			Description: description,
		})

		log.Debug("Command added to scope",
			"command", command,
			"scope", scopeString,
			"total_commands_in_scope", len(this.Commands[scopeString]),
		)
	}

	log.Debug("Bot command processing completed",
		"command", command,
		"total_scopes_processed", len(scopes),
	)
}

func (this *MyBotCommands) Register(bot *tgbotapi.BotAPI) error {
	totalScopes := len(this.Commands)
	log.Info("Starting bot commands registration",
		"total_scopes", totalScopes,
	)

	scopeIndex := 0
	for scope, botCommands := range this.Commands {
		scopeIndex++
		scopeSplit := strings.SplitN(scope, ":", 3)

		if len(scopeSplit) != 3 {
			log.Error("Invalid scope format",
				"scope", scope,
				"expected_parts", 3,
				"actual_parts", len(scopeSplit),
			)
			continue
		}

		scopeType := scopeSplit[0]
		scopeChatID, chatIDErr := strconv.ParseInt(scopeSplit[1], 10, 64)
		scopeUserID, userIDErr := strconv.ParseInt(scopeSplit[2], 10, 64)

		if chatIDErr != nil {
			log.Warn("Failed to parse chat ID in scope",
				"scope", scope,
				"chat_id_string", scopeSplit[1],
				"error", chatIDErr,
			)
		}
		if userIDErr != nil {
			log.Warn("Failed to parse user ID in scope",
				"scope", scope,
				"user_id_string", scopeSplit[2],
				"error", userIDErr,
			)
		}

		log.Debug("Registering commands for scope",
			"scope_index", scopeIndex,
			"total_scopes", totalScopes,
			"scope_type", scopeType,
			"chat_id", scopeChatID,
			"user_id", scopeUserID,
			"commands_count", len(botCommands),
		)

		setMyCommandsConfig := tgbotapi.NewSetMyCommands(botCommands...)
		setMyCommandsConfig.Scope = &tgbotapi.BotCommandScope{
			Type:   scopeType,
			ChatID: scopeChatID,
			UserID: scopeUserID,
		}

		// Log individual commands being registered
		for i, cmd := range botCommands {
			log.Debug("Command details",
				"scope_index", scopeIndex,
				"command_index", i+1,
				"command", cmd.Command,
				"description", cmd.Description,
			)
		}

		_, err := bot.Request(setMyCommandsConfig)
		if err != nil {
			log.Error("Failed to register commands for scope",
				"scope_index", scopeIndex,
				"scope_type", scopeType,
				"chat_id", scopeChatID,
				"user_id", scopeUserID,
				"commands_count", len(botCommands),
				"error", err,
			)
			return err
		}

		log.Debug("Commands registered successfully for scope",
			"scope_index", scopeIndex,
			"scope_type", scopeType,
			"chat_id", scopeChatID,
			"user_id", scopeUserID,
			"commands_count", len(botCommands),
		)
	}

	log.Info("Bot commands registration completed successfully",
		"total_scopes", totalScopes,
	)
	return nil
}
