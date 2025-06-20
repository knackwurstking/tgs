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

const (
	CBDataSingUpRequest = "Please, sign me up!"
)

type Data struct {
	Targets *tgs.Targets `yaml:"targets,omitempty"`
	Scopes  []tgs.Scope  `yaml:"scopes,omitempty"`
}

type PGVis struct {
	*tgbotapi.BotAPI

	data      *Data
	callbacks tgs.ReplyCallbacks
}

func New(api *tgbotapi.BotAPI) *PGVis {
	return &PGVis{
		data: &Data{
			Targets: tgs.NewTargets(),
			Scopes:  make([]tgs.Scope, 0),
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
	mbc.Add("/pgvissingup", "Get an api key for the \"PG Vis Server\" project.", p.data.Scopes)
}

func (p *PGVis) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		if update.CallbackQuery == nil {
			return false
		}

		slog.Debug("Got a CallbackQuery",
			"extension", p.Name(),
			"CallbackQuery.Data", update.CallbackQuery.Data,
		)

		switch update.CallbackQuery.Data {
		case CBDataSingUpRequest:
			return true
		default:
			return false
		}
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

	if update.CallbackQuery != nil {
		switch queryData := update.CallbackQuery.Data; queryData {
		case CBDataSingUpRequest:
			slog.Debug("@TODO: Need to check targets here",
				"ChatInstance", update.CallbackQuery.ChatInstance,
				"From.ID", update.CallbackQuery.From.ID,
				"From.UserName", update.CallbackQuery.From.UserName,
				"From.FirstName", update.CallbackQuery.From.FirstName,
				"From.LastName", update.CallbackQuery.From.LastName,
				"Message.Chat.ID", update.CallbackQuery.Message.Chat.ID,
			)

			if ok := tgs.CheckCallbackQueryTargets(update.CallbackQuery, p.data.Targets); !ok {
				return errors.New("invalid target")
			}

			// TODO: Get user (ID, Name), send a private
			// 		 message with the api key to the user (From)
		default:
			return fmt.Errorf("unknown callback query data: %s", queryData)
		}

		return nil
	}

	if message != nil {
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
			cbData := CBDataSingUpRequest
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

	return fmt.Errorf("there is nothing to do here")
}
