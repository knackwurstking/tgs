package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/tgs/internal/commands"
	"github.com/knackwurstking/tgs/pkg/data"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

const (
	BotCommandIP          string = "/ip"
	BotCommandJournalList string = "/journallist"
	BotCommandJournal     string = "/journal"
	BotCommandPicowStatus string = "/picowstatus"
	BotCommandPicowON     string = "/picowon"
	BotCommandPicowOFF    string = "/picowoff"
	BotCommandOPManga     string = "/opmanga"
	BotCommandOPMangaList string = "/opmangalist"
)

var (
	handledIDs = make([]int, 0) // Contains update ids already handled
)

func main() {
	app := cli.App{
		Name:  "tgs-server",
		Usage: cli.Usage("Telegram scripts server, the scripts part was kicked from the project :)"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			var configPath string

			cli.StringVar(cmd, &configPath, "config",
				cli.Usage("Path to server configuration (yaml)"),
				cli.Required,
			)

			return func(cmd *cli.Command) error {
				config := NewConfig()

				if err := loadConfig(config, configPath); err != nil {
					return err
				}

				if err := checkConfig(config); err != nil {
					return err
				}

				requestTimeout := 60 // 1 Minute
				getUpdates := tgs.RequestGetUpdates{
					API:     tgs.NewTelegramBotAPI(config.Token),
					Timeout: &requestTimeout,
				}
				for {
					resp, err := getUpdates.Send()
					if err != nil {
						slog.Warn("Request updates", "error", err)
						continue
					}

					if !resp.OK {
						slog.Error("Request updates", "response", *resp)
						return fmt.Errorf("request updates failed")
					}

					handleUpdates(config, resp.Result)
				}
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.0.0"),
		},
	}

	app.HandleError(app.Run())
}

func handleUpdates(config *Config, result []data.Update) {
	defer func() {
		newHandledIDs := make([]int, 0)
		for _, handledID := range handledIDs {
			for _, update := range result {
				if update.UpdateID == handledID {
					newHandledIDs = append(newHandledIDs, handledID)
					break
				}
			}
		}
		handledIDs = newHandledIDs
	}()

	for _, update := range result {
		if !isNewUpdateID(update.UpdateID) {
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.Entities == nil || update.Message.Text == "" {
			continue
		}

		botCommand := ""
		for _, entity := range update.Message.Entities {
			if entity.Type == "bot_command" {
				botCommand = update.Message.Text[entity.Offset:entity.Length]
				break
			}
		}

		commandConfig, err := config.Commands.Get(botCommand)
		if err != nil {
			slog.Warn("Command not found", "command", botCommand, "error", err)
			continue
		}

		if !isValidTarget(*update.Message, commandConfig.Targets, config.Targets) {
			continue
		}

		slog.Debug("Handle bot command", "command", botCommand, "message", *update.Message)

		// TODO: Handle commands
		switch botCommand {
		case BotCommandIP:
			ip := commands.NewIP(tgs.NewTelegramBotAPI(config.Token), nil)
			ip.Run(update.Message.Chat.ID)
			break
		}
	}
}

func isValidTarget(message data.Message, commandTargets *Targets, fallbackTargets Targets) bool {
	targets := fallbackTargets
	if commandTargets != nil {
		targets = *commandTargets
	}

	if message.From.ID != 0 && targets.Users != nil {
		for _, user := range targets.Users {
			if user.ID == message.From.ID {
				return true
			}
		}
	}

	if targets.Chats != nil {
		for _, chat := range targets.Chats {
			if chat.ID == message.Chat.ID {
				return true
			}
		}
	}

	return false
}

func isNewUpdateID(id int) bool {
	for _, handledID := range handledIDs {
		if handledID == id {
			return false
		}
	}

	return true
}

func checkConfig(config *Config) error {
	if config.Token == "" {
		return fmt.Errorf("token missing")
	}

	if config.Targets.Chats == nil && config.Targets.Users == nil {
		return fmt.Errorf("missing targets")
	}

	return nil
}

func loadConfig(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, config)
}
