package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
				cfg := config.New()

				if err := loadConfig(cfg, configPath); err != nil {
					return err
				}

				if err := checkConfig(cfg); err != nil {
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

				bot, err := tgbotapi.NewBotAPI(cfg.Token)
				if err != nil {
					return err
				}

				bot.Debug = false
				slog.Info("Authorized bot", "username", bot.Self.UserName)

				if err := registerCommands(bot, cfg); err != nil {
					return err
				}

				updateConfig := tgbotapi.NewUpdate(0)
				updateConfig.Timeout = 60 // 1min

				for update := range bot.GetUpdatesChan(updateConfig) {
					updateConfig.Offset = update.UpdateID + 1

					if !update.Message.IsCommand() {
						continue
					}

					// TODO: Add commands here...
					switch update.Message.Command() {
					case config.BotCommandIP[1:]:
						runCommand(
							update.Message,
							config.BotCommandIP,
							botcommand.NewIP(bot),
							cfg.IP.ValidationTargets,
						)
						break

					case config.BotCommandStats[1:]:
						runCommand(
							update.Message,
							config.BotCommandIP,
							botcommand.NewStats(bot),
							cfg.Stats.ValidationTargets,
						)
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

func runCommand(
	message *tgbotapi.Message,
	command string,
	handler botcommand.Handler,
	targets *config.ValidationTargets,
) {
	if !isValidTarget(message, targets) {
		return
	}

	logCommand(command, message)

	if err := handler.Run(message); err != nil {
		slog.Error("Command failed!", "command", command, "error", err)
	}
}

func registerCommands(bot *tgbotapi.BotAPI, cfg *config.Config) error {
	myBotCommands := tgs.NewMyBotCommands()

	myBotCommands.Add(
		config.BotCommandIP, "Get server IP",
		cfg.IP.Register,
	)

	myBotCommands.Add(
		config.BotCommandStats, "Get ID info",
		cfg.Stats.Register,
	)

	return myBotCommands.Register(bot)
}

func logCommand(command string, message *tgbotapi.Message) {
	slog.Debug("Running command.",
		"command", command,
		"message.from.username", message.From.UserName,
		"message.from.id", message.From.ID,
		"message.chat.id", message.Chat.ID,
		"message.chat.title", message.Chat.Title,
		"message.chat.type", message.Chat.Type,
		"message.message_thread_id", message.MessageThreadID,
	)
}

func isValidTarget(message *tgbotapi.Message, targets *config.ValidationTargets) bool {
	if targets.All {
		return true
	}

	if message.From.ID != 0 && targets.Users != nil {
		for _, user := range targets.Users {
			if user.ID == message.From.ID {
				return true
			}
		}
	}

	if targets.Chats != nil {
		for _, targetChat := range targets.Chats {
			if targetChat.ID == message.Chat.ID &&
				(targetChat.Type == message.Chat.Type && targetChat.Type != "") {

				if message.Chat.IsForum &&
					(message.MessageThreadID != targetChat.MessageThreadID &&
						targetChat.MessageThreadID > 0) {

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
