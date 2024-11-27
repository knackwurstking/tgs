package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SuperPaintman/nice/cli"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	botcommands "github.com/knackwurstking/tgs/internal/bot-commands"
	mybotcommands "github.com/knackwurstking/tgs/pkg/my-bot-commands"
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

				bot, err := tgbotapi.NewBotAPI(config.Token)
				if err != nil {
					return err
				}

				slog.Info("Authorized bot", "username", bot.Self.UserName)

				myBotCommands := mybotcommands.New()

				myBotCommands.Add(BotCommandIP, "Get server ip", []tgbotapi.BotCommandScope{
					{Type: "chat", ChatID: -1002493320266},
				})

				if err := myBotCommands.Register(bot); err != nil {
					return err
				}

				update := tgbotapi.NewUpdate(0)
				update.Timeout = 60 // 1min

				for update := range bot.GetUpdatesChan(update) {
					if !update.Message.IsCommand() {
						continue
					}

					switch update.Message.Command() {
					case BotCommandIP:
						if isValidTarget(update.Message, config.IPCommandConfig.Targets) {
							continue
						}

						if err := botcommands.NewIP(bot).Run(
							update.Message.Chat.ID, &update.Message.MessageID,
						); err != nil {
							slog.Error("Command failed", "command", BotCommandIP, "error", err)
						}

						break

					case BotCommandJournalList:
					case BotCommandJournal:

					case BotCommandPicowStatus:
					case BotCommandPicowON:
					case BotCommandPicowOFF:

					case BotCommandOPMangaList:
					case BotCommandOPManga:

					default:
						slog.Warn("Command not found", "command", update.Message.Command())
					}
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

func isValidTarget(message *tgbotapi.Message, targets *TargetsConfig) bool {
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
