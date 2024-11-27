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
	"github.com/knackwurstking/tgs/internal/config"
	mybotcommands "github.com/knackwurstking/tgs/pkg/my-bot-commands"
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
				cfg := config.New()

				if err := loadConfig(cfg, configPath); err != nil {
					return err
				}

				if err := checkConfig(cfg); err != nil {
					return err
				}

				slog.SetDefault(
					slog.New(
						slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
							AddSource: true,
							Level:     slog.LevelDebug,
						}),
					),
				)

				bot, err := tgbotapi.NewBotAPI(cfg.Token)
				if err != nil {
					return err
				}

				bot.Debug = true
				slog.Info("Authorized bot", "username", bot.Self.UserName)

				myBotCommands := mybotcommands.New()

				myBotCommands.Add(
					config.BotCommandIP, "Get server ip",
					cfg.IPCommandConfig.Register,
				)

				if err := myBotCommands.Register(bot); err != nil {
					return err
				}

				updateConfig := tgbotapi.NewUpdate(0)
				updateConfig.Timeout = 60 // 1min

				for update := range bot.GetUpdatesChan(updateConfig) {
					updateConfig.Offset = update.UpdateID + 1

					if !update.Message.IsCommand() {
						continue
					}

					switch update.Message.Command() {
					case config.BotCommandIP:
						if isValidTarget(update.Message, cfg.IPCommandConfig.ValidationsConfig) {
							continue
						}

						if err := botcommands.NewIP(bot).Run(
							update.Message.Chat.ID, &update.Message.MessageID,
						); err != nil {
							slog.Error("Command failed", "command", config.BotCommandIP, "error", err)
						}

						break

					case config.BotCommandJournalList:
						// TODO: ...
						break

					case config.BotCommandJournal:
						// TODO: ...
						break

					case config.BotCommandPicowStatus:
						// TODO: ...
						break

					case config.BotCommandPicowON:
						// TODO: ...
						break

					case config.BotCommandPicowOFF:
						// TODO: ...
						break

					case config.BotCommandOPMangaList:
						// TODO: ...
						break

					case config.BotCommandOPManga:
						// TODO: ...
						break

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

func isValidTarget(message *tgbotapi.Message, targets *config.ValidationsConfig) bool {
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

func checkConfig(cfg *config.Config) error {
	if cfg.Token == "" {
		return fmt.Errorf("missing token")
	}

	return nil
}

func loadConfig(cfg *config.Config, path string) error {
	extension := filepath.Ext(path)
	if extension != ".yaml" && extension != ".json" {
		return fmt.Errorf("unknown file type: %s", extension)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if extension == ".yaml" {
		return yaml.Unmarshal(data, cfg)
	}

	return json.Unmarshal(data, cfg)
}
