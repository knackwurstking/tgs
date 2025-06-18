package pgvis

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Data struct {
	Targets  *tgs.Targets          `yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `yaml:"register,omitempty"`
}

type PGVis struct {
	*tgbotapi.BotAPI

	data      *Data
	callbacks tgs.ReplyCallbacks
}

func New(api *tgbotapi.BotAPI) *PGVis {
	return &PGVis{
		data: &Data{
			Targets:  tgs.NewTargets(),
			Register: make([]tgs.BotCommandScope, 0),
		},
		callbacks: tgs.ReplyCallbacks{},
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (p *PGVis) Name() string {
	return "pgvis"
}

func (p *PGVis) SetBot(api *tgbotapi.BotAPI) {
	p.BotAPI = api
}

func (p *PGVis) ConfigPath() string {
	return "pgvis.yaml"
}

func (p *PGVis) MarshalYAML() (any, error) {
	return p.data, nil
}

func (p *PGVis) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(p.data)
}

func (p *PGVis) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/pgvissingup", "Get an api key for the \"PG Vis Server\" project.", p.data.Register)
}

func (p *PGVis) Is(message *tgbotapi.Message) bool {
	if message.ReplyToMessage != nil {
		if _, ok := p.callbacks.Get(message.ReplyToMessage.MessageID); ok {
			return true
		}
	}

	return strings.HasPrefix(message.Command(), "pgvis")
}

func (p *PGVis) Handle(message *tgbotapi.Message) error {
	if p.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	if ok := tgs.CheckTargets(message, p.data.Targets); !ok {
		return errors.New("invalid target")
	}

	command := message.Command()
	if command != "" {
	}

	switch command {
	case "pgvissingup":
		msgConfig := tgbotapi.NewMessage(message.Chat.ID, "PG Vis Server registration")
		msgConfig.ReplyToMessageID = message.MessageID

		// TODO: How to pass telegram user credentials
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonURL(
					"Sing Up", "url-to-pgvis-server-telegram-registration",
				),
			},
		)

		if _, err := p.Send(msgConfig); err != nil {
			slog.Error("Sending message failed", "extension", p.Name(), "error", err)
		}

		// TODO: Send register inline button to target(s)
		// 	- inline keyboard button
		// 	- callback url to pgvis server sing up
		// 	- pass telegram user id
		// 	- chat id
		// 	- thread id
		// 	- user name (can be empty)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}
