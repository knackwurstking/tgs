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
		startTime := time.Now()
		log.Info("Starting TGS server", "version", "2.1.1", "pid", os.Getpid())

		configHome, err := os.UserConfigDir()
		if err != nil {
			log.Error("Failed to get user config directory", "error", err)
			return err
		}

		apiConfigPath := filepath.Join(configHome, applicationName, "api.yaml")
		log.Info("Configuration setup", "config_home", configHome, "api_config", apiConfigPath)

		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.SetTimeFormat(time.RFC3339)

		c := NewConfig()

		log.Debug("Loading API configuration", "path", apiConfigPath)
		if err = loadConfig(apiConfigPath, c); err != nil {
			log.Error("Failed to load API configuration", "path", apiConfigPath, "error", err)
			return err
		}

		if c.Token == "" {
			log.Error("Missing bot token in configuration")
			return fmt.Errorf("missing token")
		}

		log.Debug("Initializing Telegram Bot API client")
		bot, err := tgbotapi.NewBotAPI(c.Token)
		if err != nil {
			log.Error("Failed to initialize Telegram Bot API", "error", err)
			return err
		}

		log.Info("Bot authenticated successfully",
			"username", bot.Self.UserName,
			"bot_id", bot.Self.ID,
			"can_join_groups", bot.Self.CanJoinGroups,
			"can_read_messages", bot.Self.CanReadAllGroupMessages,
		)
		bot.Debug = false

		// Setup extensions
		log.Info("Initializing extensions", "total_extensions", len(extensions.Register))
		for i, e := range extensions.Register {
			extensionStart := time.Now()
			if e.ConfigPath() == "" {
				log.Debug("Extension has no configuration file",
					"extension", e.Name(),
					"index", i+1,
					"total", len(extensions.Register),
				)
				continue
			}

			configPath := filepath.Join(configHome, applicationName, e.ConfigPath())
			log.Debug("Loading extension configuration",
				"extension", e.Name(),
				"config_path", configPath,
				"index", i+1,
				"total", len(extensions.Register),
			)

			err = loadConfig(configPath, e)
			if err != nil {
				log.Error("Failed to load extension configuration",
					"extension", e.Name(),
					"config_path", configPath,
					"error", err,
				)
				return err
			}

			e.SetBot(bot)
			log.Debug("Extension initialized",
				"extension", e.Name(),
				"duration_ms", time.Since(extensionStart).Milliseconds(),
			)
		}

		// Register bot commands
		log.Info("Registering bot commands")
		myBotCommands := tgs.NewMyBotCommands()

		// Add commands from extensions
		totalCommands := 0
		for _, e := range extensions.Register {
			commandsBefore := len(myBotCommands.Commands)
			log.Debug("Adding commands from extension", "extension", e.Name())
			e.AddBotCommands(myBotCommands)
			commandsAdded := len(myBotCommands.Commands) - commandsBefore
			totalCommands += commandsAdded
			log.Debug("Extension commands added",
				"extension", e.Name(),
				"commands_added", commandsAdded,
			)
		}

		log.Info("Registering commands with Telegram", "total_commands", totalCommands)
		if err = myBotCommands.Register(bot); err != nil {
			log.Error("Failed to register bot commands", "error", err)
			return err
		}
		log.Info("Bot commands registered successfully")

		// Enter the main loop
		log.Info("Server startup completed",
			"startup_duration_ms", time.Since(startTime).Milliseconds(),
		)
		log.Info("Starting update processing loop")

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60 // 1min

		updateChan := bot.GetUpdatesChan(updateConfig)
		updateCount := int64(0)

		for update := range updateChan {
			updateConfig.Offset = update.UpdateID + 1
			updateCount++
			if updateCount%100 == 0 {
				log.Info("Update processing statistics",
					"updates_processed", updateCount,
					"current_offset", updateConfig.Offset,
				)
			}
			go handleUpdate(update)
		}

		log.Info("Update processing loop ended", "total_updates_processed", updateCount)
		return nil
	}
}

func handleUpdate(update tgbotapi.Update) {
	updateStart := time.Now()
	updateID := update.UpdateID

	if update.Message != nil {
		msg := update.Message
		log.Debug("Processing message update",
			"update_id", updateID,
			"user_id", msg.From.ID,
			"username", msg.From.UserName,
			"chat_id", msg.Chat.ID,
			"chat_type", msg.Chat.Type,
			"command", msg.Command(),
			"message_id", msg.MessageID,
			"text_length", len(msg.Text),
			"has_entities", len(msg.Entities) > 0,
		)

		if msg.Chat.IsForum {
			log.Debug("Forum message details",
				"message_thread_id", msg.MessageThreadID,
			)
		}
	} else if update.CallbackQuery != nil {
		query := update.CallbackQuery
		log.Debug("Processing callback query update",
			"update_id", updateID,
			"user_id", query.From.ID,
			"username", query.From.UserName,
			"callback_data", query.Data,
			"inline_message_id", query.InlineMessageID,
		)
	} else {
		log.Debug("Processing unknown update type",
			"update_id", updateID,
			"has_message", update.Message != nil,
			"has_callback_query", update.CallbackQuery != nil,
			"has_inline_query", update.InlineQuery != nil,
			"has_channel_post", update.ChannelPost != nil,
		)
		return
	}

	handledCount := 0
	for _, e := range extensions.Register {
		if e.Is(update) {
			handledCount++
			go func(extension tgs.Extension) {
				extensionStart := time.Now()
				log.Debug("Extension processing update",
					"extension", extension.Name(),
					"update_id", updateID,
				)

				if err := extension.Handle(update); err != nil {
					log.Warn("Extension failed to handle update",
						"extension", extension.Name(),
						"update_id", updateID,
						"error", err,
						"duration_ms", time.Since(extensionStart).Milliseconds(),
					)
				} else {
					log.Debug("Extension processed update successfully",
						"extension", extension.Name(),
						"update_id", updateID,
						"duration_ms", time.Since(extensionStart).Milliseconds(),
					)
				}
			}(e)
		}
	}

	log.Debug("Update processing completed",
		"update_id", updateID,
		"extensions_handled", handledCount,
		"total_duration_ms", time.Since(updateStart).Milliseconds(),
	)
}

func loadConfig(path string, v any) error {
	loadStart := time.Now()
	log.Debug("Loading configuration file", "path", path)

	extension := filepath.Ext(path)
	if extension != ".yaml" {
		log.Error("Unsupported configuration file type",
			"path", path,
			"extension", extension,
			"supported", "yaml",
		)
		return fmt.Errorf("unknown file type: %s", extension)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error("Failed to read configuration file",
			"path", path,
			"error", err,
		)
		return err
	}

	err = yaml.Unmarshal(data, v)
	if err != nil {
		log.Error("Failed to parse YAML configuration",
			"path", path,
			"error", err,
			"file_size", len(data),
		)
		return err
	}

	log.Debug("Configuration loaded successfully",
		"path", path,
		"file_size", len(data),
		"duration_ms", time.Since(loadStart).Milliseconds(),
	)
	return nil
}
