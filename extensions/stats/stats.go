package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Data struct {
	Targets *tgs.Targets `yaml:"targets,omitempty"`
	Scopes  []tgs.Scope  `yaml:"scopes,omitempty"`
}

type Stats struct {
	*tgbotapi.BotAPI

	data *Data
}

func New(api *tgbotapi.BotAPI) *Stats {
	return &Stats{
		BotAPI: api,
		data: &Data{
			Targets: tgs.NewTargets(),
			Scopes:  make([]tgs.Scope, 0),
		},
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (s *Stats) Name() string {
	return "Stats"
}

func (s *Stats) SetBot(api *tgbotapi.BotAPI) {
	s.BotAPI = api
}

func (s *Stats) ConfigPath() string {
	return "stats.yaml"
}

func (s *Stats) MarshalYAML() (any, error) {
	return s.data, nil
}

func (s *Stats) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(s.data)
}

func (s *Stats) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/stats", "Get ID info", s.data.Scopes)
}

func (s *Stats) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	return strings.HasPrefix(update.Message.Command(), "stats")
}

func (s *Stats) Handle(update tgbotapi.Update) error {
	if s.BotAPI == nil {
		log.Error("Stats extension BotAPI is nil")
		panic("BotAPI is nil!")
	}

	message := update.Message
	log.Debug("Stats extension handling update",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
		"command", message.Command(),
	)

	if ok := tgs.CheckTargets(message, s.data.Targets); !ok {
		log.Debug("Stats request from unauthorized target",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
			"command", message.Command(),
		)
		return errors.New("invalid target")
	}

	if command := message.Command(); command != "stats" {
		log.Warn("Unknown command in Stats extension",
			"command", command,
			"user_id", message.From.ID,
		)
		return fmt.Errorf("unknown command: %s", command)
	}

	log.Info("Processing stats request",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
		"message_thread_id", message.MessageThreadID,
	)

	data := struct {
		UserName        string `json:"username"`
		UserID          int64  `json:"user_id"`
		ChatID          int64  `json:"chat_id"`
		MessageThreadID int    `json:"message_thread_id"`
	}{
		UserName:        message.From.UserName,
		UserID:          message.From.ID,
		ChatID:          message.Chat.ID,
		MessageThreadID: message.MessageThreadID,
	}

	log.Debug("Preparing stats data response",
		"username", data.UserName,
		"user_id", data.UserID,
		"chat_id", data.ChatID,
		"message_thread_id", data.MessageThreadID,
	)

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error("Failed to marshal stats data to JSON",
			"error", err,
			"user_id", message.From.ID,
		)
		return err
	}

	msgConfig := tgbotapi.NewMessage(message.Chat.ID,
		"```json\n"+fmt.Sprintf("%s\n", string(jsonData))+"```",
	)

	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	log.Debug("Sending stats response",
		"chat_id", message.Chat.ID,
		"reply_to_message_id", message.MessageID,
		"json_size_bytes", len(jsonData),
	)

	_, err = s.Send(msgConfig)
	if err != nil {
		log.Error("Failed to send stats response",
			"chat_id", message.Chat.ID,
			"user_id", message.From.ID,
			"error", err,
		)
	} else {
		log.Info("Stats response sent successfully",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
		)
	}
	return err
}
