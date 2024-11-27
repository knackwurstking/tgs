package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SuperPaintman/nice/cli"
	"github.com/knackwurstking/tgs/internal/commands"
	"github.com/knackwurstking/tgs/pkg/data"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
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

				if err := setBotCommands(config); err != nil {
					return err
				}

				if err := updateLoop(config); err != nil {
					return err
				}

				return nil
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.0.0"),
		},
	}

	app.HandleError(app.Run())
}

func updateLoop(config *Config) error {
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

// TODO: Need to clear all existing commands first
func setBotCommands(config *Config) error {
	api := tgs.NewTelegramBotAPI(config.Token)
	requests := make([]tgs.RequestSetMyCommands, 0)

	getRequest := func(scope Scope) *tgs.RequestSetMyCommands {
		for i, r := range requests {
			if r.Scope.Type == scope.Scope &&
				r.Scope.UserID == scope.UserID &&
				r.Scope.ChatID == scope.ChatID {

				return &requests[i]
			}
		}

		return nil
	}

	addCommandToRequests := func(scopes []Scope, botCommand data.BotCommand) {
		if scopes == nil {
			return
		}

		var request *tgs.RequestSetMyCommands

		for _, scope := range config.Commands.IP.Scopes {
			if request = getRequest(scope); request != nil {
				request.Commands = append(request.Commands, botCommand)
			} else {
				requests = append(
					requests,
					tgs.RequestSetMyCommands{
						API:      api,
						Commands: []data.BotCommand{botCommand},
						Scope: data.BotCommandScope{
							Type:   scope.Scope,
							ChatID: scope.ChatID,
							UserID: scope.UserID,
						},
					},
				)
			}
		}
	}

	if !config.Commands.IP.Disabled {
		addCommandToRequests(config.Commands.IP.Scopes, data.BotCommand{
			Command:     BotCommandIP,
			Description: "Get server ip",
		})
	}

	if !config.Commands.JournalList.Disabled {
		addCommandToRequests(config.Commands.IP.Scopes, data.BotCommand{
			Command:     BotCommandJournalList,
			Description: "List available journals",
		})
	}

	if !config.Commands.Journal.Disabled {
		addCommandToRequests(config.Commands.Journal.Scopes, data.BotCommand{
			Command:     BotCommandJournal,
			Description: "Get a journal",
		})
	}

	if !config.Commands.PicowStatus.Disabled {
		addCommandToRequests(config.Commands.PicowStatus.Scopes, data.BotCommand{
			Command:     BotCommandPicowStatus,
			Description: "Lights power status",
		})
	}

	if !config.Commands.PicowOn.Disabled {
		addCommandToRequests(config.Commands.PicowOn.Scopes, data.BotCommand{
			Command:     BotCommandPicowON,
			Description: "Power on the lights",
		})
	}

	if !config.Commands.PicowOff.Disabled {
		addCommandToRequests(config.Commands.PicowOff.Scopes, data.BotCommand{
			Command:     BotCommandPicowOFF,
			Description: "Power off the lights",
		})
	}

	if !config.Commands.OPMangaList.Disabled {
		addCommandToRequests(config.Commands.OPMangaList.Scopes, data.BotCommand{
			Command:     BotCommandOPMangaList,
			Description: "List One Piece mangas available",
		})
	}

	if !config.Commands.OPManga.Disabled {
		addCommandToRequests(config.Commands.OPManga.Scopes, data.BotCommand{
			Command:     BotCommandOPManga,
			Description: "Get a One Piece manga episode",
		})
	}

	for _, request := range requests {
		resp, err := request.Send()
		if err != nil {
			return err
		}
		if !resp.OK || !resp.Result {
			return fmt.Errorf("Set commands failed on Telegram %+v", resp)
		}
	}

	return fmt.Errorf("under construction")
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

		if !isValidTarget(*update.Message, commandConfig.Targets) {
			continue
		}

		slog.Debug("Handle bot command", "command", botCommand, "message", *update.Message)

		var tgCommandHandler commands.TelegramCommandHandler
		switch botCommand {
		case BotCommandIP:
			tgCommandHandler = commands.NewIP(tgs.NewTelegramBotAPI(config.Token))
			break

		case BotCommandJournalList:
			// TODO: ...
			break
		case BotCommandJournal:
			// TODO: ...
			break

		case BotCommandPicowStatus:
			// TODO: ...
			break
		case BotCommandPicowON:
			// TODO: ...
			break
		case BotCommandPicowOFF:
			// TODO: ...
			break

		case BotCommandOPMangaList:
			// TODO: ...
			break
		case BotCommandOPManga:
			// TODO: ...
			break
		}

		if tgCommandHandler != nil {
			tgCommandHandler.Run(update.Message.Chat.ID)
		}
	}
}

func isValidTarget(message data.Message, targets *Targets) bool {
	if message.From.ID != 0 && targets.Users != nil {
		for _, user := range targets.Users {
			if user.ID == message.From.ID {
				return true
			}
		}
	}

	if targets.Chats != nil {
		for _, chat := range targets.Chats {
			if chat.ID == message.Chat.ID && (chat.Type == message.Chat.Type && chat.Type != "") {
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
		return fmt.Errorf("missing token")
	}

	return nil
}

func loadConfig(config *Config, path string) error {
	extension := filepath.Ext(path)
	if extension != ".yaml" && extension != ".json" {
		return fmt.Errorf("unknown file type: %s", extension)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if extension == ".yaml" {
		return yaml.Unmarshal(data, config)
	}

	return json.Unmarshal(data, config)
}
