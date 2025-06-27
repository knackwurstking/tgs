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
			slog.Debug("Callback Query",
				"Data", update.CallbackQuery.Data,
				"From.ID", update.CallbackQuery.From.ID,
				"From.UserName", update.CallbackQuery.From.UserName,
				"From.FirstName", update.CallbackQuery.From.FirstName,
				"From.LastName", update.CallbackQuery.From.LastName,
				"Message.Chat.ID", update.CallbackQuery.Message.Chat.ID,
			)

			if ok := tgs.CheckCallbackQueryTargets(update.CallbackQuery, p.data.Targets); !ok {
				return errors.New("invalid target")
			}

			user, err := NewUser(
				update.CallbackQuery.From.ID,
				update.CallbackQuery.From.UserName,
			)
			if err != nil {
				msgConfig := tgbotapi.NewMessage(
					update.CallbackQuery.From.ID,
					fmt.Sprintf(
						"Ups, etwas ist schief gelaufen: \n`%s`",
						err.Error(),
					),
				)
				msgConfig.ParseMode = "MarkdownV2"

				if _, err = p.Send(msgConfig); err != nil {
					slog.Error("Sending message failed", "extension", p.Name(), "error", err)
				}

				return nil
			}

			// Into message
			msgConfig := tgbotapi.NewMessage(
				update.CallbackQuery.From.ID,
				fmt.Sprintf(
					"Den folgenden Api Key bitte sicher aufbewahren, "+
						"das ist dein zugang zur App. Dieser key ist gebunden "+
						"an deine Telegram ID.",
				),
			)

			if _, err := p.Send(msgConfig); err != nil {
				slog.Error("Sending message failed", "extension", p.Name(), "error", err)
			}

			// Api key message
			msgConfig = tgbotapi.NewMessage(update.CallbackQuery.From.ID, user.ApiKey)

			if _, err := p.Send(msgConfig); err != nil {
				slog.Error("Sending message failed", "extension", p.Name(), "error", err)
			}

			// Link the sing up page now
			msgConfig = tgbotapi.NewMessage(update.CallbackQuery.From.ID,
				"Zur registrierung gehts hier lang")

			msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonURL(
						"Registrierung",
						fmt.Sprintf(
							"https://knackwurstking.com/pg-vis/signup?key=%s",
							user.ApiKey,
						),
					),
				},
			)

			if _, err := p.Send(msgConfig); err != nil {
				slog.Error("Sending message failed", "extension", p.Name(), "error", err)
			}
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
			msgConfig := tgbotapi.NewMessage(
				message.Chat.ID,
				fmt.Sprintf(
					"Wenn du auf den \"Sing Up\" Button klickst, "+
						"bekommst du deinen Api Key zugesendet.\n\n"+
						"Bitte ignorieren, immer noch am testen.",
				),
			)

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
