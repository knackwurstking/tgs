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
					cfg = config.New(bot, make(chan *botcommand.Reply))
				)

				slog.SetDefault(
					slog.New(
						tint.NewHandler(os.Stderr, &tint.Options{
							AddSource: true,
							Level:     slog.LevelDebug,
						}),
					),
				)

				if err = loadConfig(cfg, configPath); err != nil {
					return err
				}

				if err = checkConfig(cfg); err != nil {
					return err
				}

				bot, err = tgbotapi.NewBotAPI(cfg.Token)
				if err != nil {
					return err
				}

				bot.Debug = false
				cfg.IP.BotAPI = bot
				cfg.Stats.BotAPI = bot
				cfg.Journal.BotAPI = bot
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

				replyCallbacks := map[int]*botcommand.Reply{}
				updateChan := bot.GetUpdatesChan(updateConfig)
				for {
					select {
					case update := <-updateChan:
						updateConfig.Offset = update.UpdateID + 1

						if !update.Message.IsCommand() {
							replyID := update.Message.ReplyToMessage.MessageID
							if r, ok := replyCallbacks[replyID]; ok {
								r.Run(update.Message)
								continue
							}

							slog.Debug("Got a new update",
								"replyID", replyID,
								"update.Message.Text", update.Message.Text,
							)
							continue
						}

						// Run commands here
						switch v := update.Message.Command(); {
						case strings.HasPrefix(v, botcommand.BotCommandIP[1:]):
							runCommand(cfg.IP, update.Message)
							break

						case strings.HasPrefix(v, botcommand.BotCommandStats[1:]):
							runCommand(cfg.Stats, update.Message)
							break

						case strings.HasPrefix(v, botcommand.BotCommandJournal[1:]):
							runCommand(cfg.Journal, update.Message)
							break

						default:
							slog.Warn("Command not found!", "command", v)
						}

						break

					case reply := <-cfg.Reply:
						messageID := reply.Message.MessageID
						slog.Debug("Set a reply callback function", "messageID", messageID)

						if r, ok := replyCallbacks[messageID]; ok {
							r.Done() <- nil
						}

						replyCallbacks[messageID] = reply
						go reply.StartTimeout()

						go func() {
							defer reply.Close()

							switch err := <-reply.Done(); {
							case err == botcommand.ReplyTimeoutError:
								slog.Warn("Reply callback timeout",
									"messageID", messageID, "reply.Timeout", reply.Timeout,
								)
								break

							case err == nil:
								slog.Debug("Reply callback finished",
									"messageID", messageID, "reply.Timeout", reply.Timeout,
								)

							default:
								slog.Warn("Reply callback finished",
									"messageID", messageID, "reply.Timeout", reply.Timeout,
									"error", err,
								)

								msgConfig := tgbotapi.NewMessage(
									reply.Message.Chat.ID,
									fmt.Sprintf("`%s`", err.Error()),
								)
								msgConfig.ReplyToMessageID = messageID
								msgConfig.ParseMode = "MarkdownV2"

								_, _ = bot.Send(msgConfig) // NOTE: Ignore any error

								break
							}

							delete(replyCallbacks, messageID)
						}()

						break
					}
				}
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.1.0"),
		},
	}

	app.HandleError(app.Run())
}

func runCommand(handler botcommand.Handler, message *tgbotapi.Message) {
	if !isValidTarget(message, handler) {
		return
	}

	command := message.Command()
	slog.Debug("Running command.",
		"command", command,
		"message.message_id", message.MessageID,
		"message.from.username", message.From.UserName,
		"message.from.id", message.From.ID,
		"message.chat.id", message.Chat.ID,
		"message.chat.title", message.Chat.Title,
		"message.chat.type", message.Chat.Type,
		"message.message_thread_id", message.MessageThreadID,
	)

	go func() {
		if err := handler.Run(message); err != nil {
			slog.Error("Command failed!", "command", command, "error", err)
		}
	}()
}

func isValidTarget(message *tgbotapi.Message, handler botcommand.Handler) bool {
	if handler.Targets() == nil {
		return false
	}

	if handler.Targets().All {
		return true
	}

	// User ID check
	if message.From.ID != 0 && handler.Targets().Users != nil {
		for _, user := range handler.Targets().Users {
			if user.ID == message.From.ID {
				return true
			}
		}
	}

	// Chat ID check & message thread ID if chat is forum
	if handler.Targets().Chats != nil {
		for _, t := range handler.Targets().Chats {
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
