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

func (p *PGVis) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		if update.CallbackQuery != nil {
			slog.Debug("Got a CallbackQuery",
				"extension", p.Name(),
				"CallbackQuery.Data", update.CallbackQuery.Data,
				"CallbackQuery.Message", update.CallbackQuery.Message,
			)
			// TODO: Validate `CallbackQuery.Data`, or whatever
			return true
		}

		return false
	}

	if update.Message.ReplyToMessage != nil {
		if _, ok := p.callbacks.Get(update.Message.ReplyToMessage.MessageID); ok {
			return true
		}
	}

	return strings.HasPrefix(update.Message.Command(), "pgvis")
}

func (p *PGVis) Handle(update tgbotapi.Update) error {
	if p.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	message := update.Message

	// TODO:Handle callback queries here
	if message == nil {
		return nil
	}

	if ok := tgs.CheckTargets(message, p.data.Targets); !ok {
		return errors.New("invalid target")
	}

	switch command := message.Command(); command {
	case "pgvissingup":
		msgConfig := tgbotapi.NewMessage(message.Chat.ID, "Bitte ignorieren, bin am testen!")
		msgConfig.ReplyToMessageID = message.MessageID

		button := tgbotapi.NewInlineKeyboardButtonData(
			"Sing Up", "/singup",
		)
		cbData := "Please, sign me up!"
		button.CallbackData = &cbData

		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				button,
			},
		)

		if _, err := p.Send(msgConfig); err != nil {
			slog.Error("Sending message failed", "extension", p.Name(), "error", err)
		}
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}
