package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SuperPaintman/nice/cli"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lmittmann/tint"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/extensions"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

const (
	applicationName = "tgs-server"
)

func main() {
	app := cli.App{
		Name:  applicationName,
		Usage: cli.Usage("Telegram scripts server, the scripts part was kicked from the project :)"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			return actionHandler()
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("2.0.0"),
		},
	}

	app.HandleError(app.Run())
}

func actionHandler() func(cmd *cli.Command) error {
	return func(cmd *cli.Command) error {
		var (
			err error
			bot *tgbotapi.BotAPI
		)

		slog.SetDefault(
			slog.New(
				tint.NewHandler(os.Stderr, &tint.Options{
					AddSource: true,
					Level:     slog.LevelDebug,
				}),
			),
		)

		var configHome string
		if configHome, err = os.UserConfigDir(); err != nil {
			return err
		}

		apiConfigPath := filepath.Join(configHome, applicationName, "api.yaml")
		slog.Debug("API Config location", "path", apiConfigPath)

		c := NewConfig()
		if err = loadConfig(apiConfigPath, c); err != nil {
			return err
		}
		if c.Token == "" {
			return fmt.Errorf("missing token")
		}

		bot, err = tgbotapi.NewBotAPI(c.Token)
		if err != nil {
			return err
		}

		slog.Info("Authorized bot", "username", bot.Self.UserName)
		bot.Debug = false

		// Setup bots
		for _, e := range extensions.Register {
			if e.ConfigPath() == "" {
				slog.Debug("Skip config for extension", "name", e.Name())
				continue
			}

			configPath := filepath.Join(configHome, applicationName, e.ConfigPath())
			slog.Debug("Try to load extension configuration",
				"name", e.Name(), "path", configPath)

			err = loadConfig(configPath, e)
			if err != nil {
				return err
			}

			e.SetBot(bot)
		}

		// Register bot commands here
		myBotCommands := tgs.NewMyBotCommands()

		// Add commands from extension
		for _, e := range extensions.Register {
			e.AddBotCommands(myBotCommands)
		}

		if err = myBotCommands.Register(bot); err != nil {
			return err
		}

		for _, e := range extensions.Register {
			e.Start() // NOTE: Maybe using a goroutine here?
		}

		// Enter the main loop
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60 // 1min

		updateChan := bot.GetUpdatesChan(updateConfig)
		for {
			select {
			case update := <-updateChan:
				updateConfig.Offset = update.UpdateID + 1
				go handleUpdate(update)
			}
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message != nil {
		slog.Debug(
			"New update",
			"message.Command", update.Message.Command,
			"message.Text", update.Message.Text,
			"message.From", update.Message.From,
			"message.From.ID", update.Message.From.ID,
			"message.Chat", update.Message.Chat,
			"message.ReplyToMessage", update.Message.ReplyToMessage,
		)
	}

	for _, e := range extensions.Register {
		if e.Is(update.Message) {
			go func() {
				slog.Info("Handle update",
					"extension", e.Name(),
					"command", update.Message.Command(),
					"message.MessageID", update.Message.MessageID,
					"message.Chat", update.Message.Chat,
					"message.From", update.Message.From,
					"message.From.ID", update.Message.From.ID,
				)

				if err := e.Handle(update.Message); err != nil {
					slog.Warn("Handle update failed",
						"extension", e.Name(),
						"message.MessageID", update.Message.MessageID,
						"error", err,
					)
				}
			}()
		}
	}
}

func loadConfig(path string, v any) error {
	extension := filepath.Ext(path)
	if extension != ".yaml" {
		return fmt.Errorf("unknown file type: %s", extension)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}
