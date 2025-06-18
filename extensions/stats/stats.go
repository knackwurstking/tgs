package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Data struct {
	Targets  *tgs.Targets          `yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `yaml:"register,omitempty"`
}

type Stats struct {
	*tgbotapi.BotAPI

	data *Data
}

func New(api *tgbotapi.BotAPI) *Stats {
	return &Stats{
		BotAPI: api,
		data: &Data{
			Targets:  tgs.NewTargets(),
			Register: make([]tgs.BotCommandScope, 0),
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
	mbc.Add("/stats", "Get ID info", s.data.Register)
}

func (s *Stats) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	return strings.HasPrefix(update.Message.Command(), "stats")
}

func (s *Stats) Handle(update tgbotapi.Update) error {
	if s.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	message := update.Message

	if ok := tgs.CheckTargets(message, s.data.Targets); !ok {
		return errors.New("invalid target")
	}

	if command := message.Command(); command != "stats" {
		return fmt.Errorf("unknown command: %s", command)
	}

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

	jsonData, err := json.MarshalIndent(data, "", "    ")

	msgConfig := tgbotapi.NewMessage(message.Chat.ID,
		"```json\n"+fmt.Sprintf("%s\n", string(jsonData))+"```",
	)

	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = s.Send(msgConfig)
	return err
}
