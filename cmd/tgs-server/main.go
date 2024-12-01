package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/SuperPaintman/nice/cli"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lmittmann/tint"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/config"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

func main() {
	app := cli.App{
		Name:  "tgs-server",
		Usage: cli.Usage("Telegram scripts server, the scripts part was kicked from the project :)"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			var configPath string

			cli.StringVar(cmd, &configPath, "config",
				cli.Usage("Path to server configuration (yaml)"),
				cli.WithShort("c"),
				cli.Required,
			)

			return func(cmd *cli.Command) error {
				var (
					err error
					bot *tgbotapi.BotAPI
					cfg = config.New(bot)
				)

				if err = loadConfig(cfg, configPath); err != nil {
					return err
				}

				if err = checkConfig(cfg); err != nil {
					return err
				}

				slog.SetDefault(
					slog.New(
						tint.NewHandler(os.Stderr, &tint.Options{
							AddSource: true,
							Level:     slog.LevelDebug,
						}),
					),
				)

				bot, err = tgbotapi.NewBotAPI(cfg.Token)
				if err != nil {
					return err
				}

				bot.Debug = false
				slog.Info("Authorized bot", "username", bot.Self.UserName)

				// Register bot commands here
				myBotCommands := tgs.NewMyBotCommands()
				cfg.IP.AddCommands(myBotCommands, cfg.IP.Register()...)
				cfg.Stats.AddCommands(myBotCommands, cfg.Stats.Register()...)
				cfg.Journal.AddCommands(myBotCommands, cfg.Journal.Register()...)

				if err = myBotCommands.Register(bot); err != nil {
					return err
				}

				// Enter the main loop
				updateConfig := tgbotapi.NewUpdate(0)
				updateConfig.Timeout = 60 // 1min

				for update := range bot.GetUpdatesChan(updateConfig) {
					updateConfig.Offset = update.UpdateID + 1

					if !update.Message.IsCommand() {
						continue
					}

					// Run commands here
					switch v := update.Message.Command(); {
					case strings.HasSuffix(v, botcommand.BotCommandIP[1:]):
						runCommand(cfg.IP, update.Message)
						break

					case strings.HasSuffix(v, botcommand.BotCommandStats[1:]):
						runCommand(cfg.Stats, update.Message)
						break

					case strings.HasSuffix(v, botcommand.BotCommandJournal[1:]):
						runCommand(cfg.Journal, update.Message)
						break

					default:
						slog.Warn("Command not found!", "command", update.Message.Command())
					}
				}

				return nil
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.2.0.dev"),
		},
	}

	app.HandleError(app.Run())
}

func runCommand(handler botcommand.Handler, message *tgbotapi.Message) {
	if !isValidTarget(message, handler.Targets()) {
		return
	}

	command := message.Command()
	slog.Debug("Running command.",
		"command", command,
		"message.from.username", message.From.UserName,
		"message.from.id", message.From.ID,
		"message.chat.id", message.Chat.ID,
		"message.chat.title", message.Chat.Title,
		"message.chat.type", message.Chat.Type,
		"message.message_thread_id", message.MessageThreadID,
	)

	if err := handler.Run(message); err != nil {
		slog.Error("Command failed!", "command", command, "error", err)
	}
}

func isValidTarget(message *tgbotapi.Message, targets *botcommand.Targets) bool {
	if targets == nil {
		return false
	}

	if targets.All {
		return true
	}

	// User ID check
	if message.From.ID != 0 && targets.Users != nil {
		for _, user := range targets.Users {
			if user.ID == message.From.ID {
				return true
			}
		}
	}

	// Chat ID check & message thread ID if chat is forum
	if targets.Chats != nil {
		for _, t := range targets.Chats {
			if t.ID == message.Chat.ID && (t.Type == message.Chat.Type && t.Type != "") {
				if message.Chat.IsForum &&
					(message.MessageThreadID != t.MessageThreadID &&
						t.MessageThreadID > 0) {

					return false
				}

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
