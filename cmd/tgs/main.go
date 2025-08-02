package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SuperPaintman/nice/cli"
	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/extensions"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

const (
	applicationName = "tgs"
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
			cli.VersionCommandFlag("2.1.1"),
		},
	}

	app.HandleError(app.Run())
}

func actionHandler() func(cmd *cli.Command) error {
	return func(cmd *cli.Command) error {
		configHome, err := os.UserConfigDir()
		if err != nil {
			return err
		}

		apiConfigPath := filepath.Join(configHome, applicationName, "api.yaml")

		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.SetTimeFormat(time.RFC3339)

		c := NewConfig()

		if err = loadConfig(apiConfigPath, c); err != nil {
			return err
		}

		if c.Token == "" {
			return fmt.Errorf("missing token")
		}

		bot, err := tgbotapi.NewBotAPI(c.Token)
		if err != nil {
			return err
		}

		log.Infof("Authorized bot with username: %s", bot.Self.UserName)
		bot.Debug = false

		// Setup bots
		for _, e := range extensions.Register {
			if e.ConfigPath() == "" {
				log.Debugf("Skip configuration for extension: %s", e.Name())
				continue
			}

			configPath := filepath.Join(configHome, applicationName, e.ConfigPath())
			log.Debugf("Try to load the configuration for the \"%s\" extension", e.Name())
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
			log.Debugf("Add bot commands for the \"%s\" extension", e.Name())
			e.AddBotCommands(myBotCommands)
		}

		if err = myBotCommands.Register(bot); err != nil {
			return err
		}

		// Enter the main loop
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60 // 1min

		updateChan := bot.GetUpdatesChan(updateConfig)
		for update := range updateChan {
			updateConfig.Offset = update.UpdateID + 1
			go handleUpdate(update)
		}

		return nil
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		log.Debugf(
			"Message: Handle update from %d: %s - %s",
			update.Message.From.ID,
			update.Message.Command(), update.Message.Text,
		)

		log.Debugf("--> From: %#v", update.Message.From)
		log.Debugf("--> Chat: %#v", update.Message.Chat)
	} else if update.CallbackQuery != nil {
		log.Debugf("CallbackQuery: Handle update from %d: %#v",
			update.Message.From.ID, update.CallbackQuery)
	} else {
		log.Debugf("Unknown: Handle update: %#v", update)
		return
	}

	for _, e := range extensions.Register {
		if e.Is(update) {
			go func() {
				log.Debugf("Extension: %s: Handle update", e.Name())

				if err := e.Handle(update); err != nil {
					log.Warnf("Extension: %s: Handle update failed: %s", e.Name(), err)
				}
			}()
		}
	}
}

func loadConfig(path string, v any) error {
	log.Debugf("Loading configuration from: %s", path)

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
