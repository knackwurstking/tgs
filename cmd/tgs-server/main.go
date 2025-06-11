// TODO: Continue here
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
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/config"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

var replyCallbacks = map[int]*botcommand.Reply{}

func main() {
	app := cli.App{
		Name:  "tgs-server",
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

		var configHome string
		if configHome, err = os.UserConfigDir(); err != nil {
			return err
		}

		if err = loadConfig(filepath.Join(configHome, "api.yaml"), cfg); err != nil {
			return err
		}

		if cfg.Token == "" {
			return fmt.Errorf("missing token")
		}

		bot, err = tgbotapi.NewBotAPI(cfg.Token)
		if err != nil {
			return err
		}

		slog.Info("Authorized bot", "username", bot.Self.UserName)
		bot.Debug = false

		// Setup bots
		for _, e := range extensions.Register {
			err = loadConfig(filepath.Join(configHome, e.ConfigPath()), e)
			if err != nil {
				return err
			}

			e.SetBot(bot)
		}

		cfg.SetBot(bot)

		// Register bot commands here
		myBotCommands := tgs.NewMyBotCommands()

		// Add commands from extension
		for _, e := range extensions.Register {
			e.Commands(myBotCommands)
		}

		// cfg.IP.AddCommands(myBotCommands)
		cfg.Stats.AddCommands(myBotCommands)
		cfg.OPManga.AddCommands(myBotCommands)

		if err = myBotCommands.Register(bot); err != nil {
			return err
		}

		// Enter the main loop
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60 // 1min

		updateChan := bot.GetUpdatesChan(updateConfig)
		for {
			select {
			case update := <-updateChan:
				updateConfig.Offset = update.UpdateID + 1
				go handleUpdate(update, cfg)

			case reply := <-cfg.Reply:
				go handleReplies(reply, bot)
			}
		}
	}
}

func handleUpdate(update tgbotapi.Update, cfg *config.Config) {
	if update.Message == nil {
		return
	}

	if len(update.Message.NewChatMembers) > 0 {
		for _, u := range update.Message.NewChatMembers {
			// Get the chat ID, I'm not sure if it'll always be set
			chatID := int64(-1)
			if update.Message.Chat != nil {
				chatID = update.Message.Chat.ID
			}

			slog.Warn("New chat member", "chat_id", chatID, "user", u)
		}
	}

	// NOTE: Handle reply callbacks (out dated)
	//
	// if !update.Message.IsCommand() {
	//	if update.Message.ReplyToMessage == nil {
	//		return
	//	}
	//
	//	replyID := update.Message.ReplyToMessage.MessageID
	//	if r, ok := replyCallbacks[replyID]; ok {
	//		r.Run(update.Message)
	//	} else {
	//		slog.Debug("Got a new update",
	//			"replyID", replyID,
	//			"update.Message.Text", update.Message.Text,
	//		)
	//	}
	//
	//	return
	//}

	for _, e := range extensions.Register {
		if e.Is(update.Message) {
			go e.Handle(update.Message)
		}
	}

	// switch {
	// case strings.HasPrefix(command, cfg.Stats.BotCommand()):
	//	handleCommand(cfg.Stats, update.Message)

	// case strings.HasPrefix(command, cfg.OPManga.BotCommand()):
	//	handleCommand(cfg.OPManga, update.Message)

	//default:
	//	slog.Warn("Command not found!", "command", command)
	//}
}

// TODO: Check for valid targets inside extension `...Is(...)`
//
//func handleCommand(handler botcommand.Handler, message *tgbotapi.Message) {
//	if !isValidTarget(message, handler) {
//		slog.Debug("Invalid target",
//			"command", message.Command(),
//			"message.Chat.ID", message.Chat.ID,
//			"handler.Targets", handler.Targets(),
//		)
//		return
//	}
//
//	command := message.Command()
//	slog.Debug("Running command.",
//		"command", command,
//		"message.message_id", message.MessageID,
//		"message.from.username", message.From.UserName,
//		"message.from.id", message.From.ID,
//		"message.chat.id", message.Chat.ID,
//		"message.chat.title", message.Chat.Title,
//		"message.chat.type", message.Chat.Type,
//		"message.message_thread_id", message.MessageThreadID,
//	)
//
//	go func() {
//		if err := handler.Run(message); err != nil {
//			slog.Error("Command failed!", "command", command, "error", err)
//		}
//	}()
//}

func handleReplies(r *botcommand.Reply, bot *tgbotapi.BotAPI) {
	messageID := r.Message.MessageID
	slog.Debug("Set a reply callback function", "messageID", messageID)

	if r, ok := replyCallbacks[messageID]; ok {
		r.Done() <- nil
	}

	replyCallbacks[messageID] = r
	go r.StartTimeout()

	defer r.Close()

	err := <-r.Done()
	switch err {
	case botcommand.ErrorTimeout:
		slog.Warn("Reply callback timeout",
			"messageID", messageID, "reply.Timeout", r.Timeout,
		)

	case nil:
		slog.Debug("Reply callback finished",
			"messageID", messageID, "reply.Timeout", r.Timeout,
		)

	default:
		slog.Warn("Reply callback finished",
			"messageID", messageID, "reply.Timeout", r.Timeout,
			"error", err,
		)

		msgConfig := tgbotapi.NewMessage(
			r.Message.Chat.ID,
			fmt.Sprintf("`%s`", err.Error()),
		)
		msgConfig.ReplyToMessageID = messageID
		msgConfig.ParseMode = "MarkdownV2"

		_, _ = bot.Send(msgConfig) // NOTE: Ignore any error
	}

	delete(replyCallbacks, messageID)
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
	chats := handler.Targets().Chats
	if chats == nil {
		return false
	}

	for _, t := range chats {
		if t.ID == message.Chat.ID && (t.Type == message.Chat.Type && t.Type != "") {
			if message.Chat.IsForum {
				if t.MessageThreadID <= 0 {
					return true
				}

				if t.MessageThreadID == message.MessageThreadID {
					return true
				}
			} else {
				return true
			}
		}
	}

	return false
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
