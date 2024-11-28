package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SuperPaintman/nice/cli"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/bot-commands/ip"
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

				bot, err := tgbotapi.NewBotAPI(cfg.Token)
				if err != nil {
					return err
				}

				bot.Debug = true
				log.Printf("Authorized bot, username=%s", bot.Self.UserName)

				myBotCommands := tgs.NewMyBotCommands()

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
					case config.BotCommandIP[1:]:
						if !isValidTarget(update.Message, cfg.IPCommandConfig.ValidationsConfig) {
							continue
						}

						log.Printf("Running command %s, from=%+v, chat=%+v",
							config.BotCommandIP, update.Message.From, update.Message.Chat,
						)

						if err := ip.New(bot).Run(
							update.Message.Chat.ID, &update.Message.MessageID,
						); err != nil {
							log.Printf("Command \"%s\" failed with: %s", config.BotCommandIP, err)
						}

						break

					case config.BotCommandStats[1:]:
						// TODO: Continue here...

						break

					default:
						log.Printf("Command \"%s\" not found!", update.Message.Command())
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

				if message.Chat.Type == "supergroup" &&
					(message.MessageThreadID != targetChat.MessageThreadID &&
						targetChat.MessageThreadID > -1) {

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
