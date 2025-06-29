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
	command := update.Message.Command()

	if command == "start" {
		return update.Message.Text == fmt.Sprintf("/start pgvissingup")
	}

	return strings.HasPrefix(command, "pgvis")
}

func (p *PGVis) Handle(update tgbotapi.Update) error {
	if p.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	message := update.Message

	if message != nil {

		switch command := message.Command(); command {
		case "start":
			if !tgs.CheckTargetsForUser(message.From.ID, p.data.Targets) {
				return errors.New("invalid target")
			}

			user, err := NewUser(message.From.ID, message.From.UserName)
			if err != nil {
				return err
			}

			// Info message
			msgConfig := tgbotapi.NewMessage(
				message.From.ID,
				fmt.Sprintf(
					"Den folgenden Api Key bitte sicher aufbewahren, "+
						"das ist dein zugang zur App. Dieser key ist gebunden "+
						"an deine Telegram ID.",
				),
			)

			if _, err := p.Send(msgConfig); err != nil {
				slog.Error("Sending message failed", "extension", p.Name(), "error", err)
			}

			// ApiKey message
			msgConfig = tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("`%s`", user.ApiKey))
			msgConfig.ParseMode = "MarkdownV2"

			if _, err := p.Send(msgConfig); err != nil {
				slog.Error("Sending message failed", "extension", p.Name(), "error", err)
			}

			// Link to the pg-vis server singup page
			msgConfig = tgbotapi.NewMessage(message.From.ID,
				"Zur registrierung gehts hier lang")

			msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonURL(
						"Registrierung",
						fmt.Sprintf(
							"https://knackwurstking.com/pg-vis/signup?key=%s", // FIXME: This link does not exists yet
							user.ApiKey,
						),
					),
				},
			)

			if _, err := p.Send(msgConfig); err != nil {
				slog.Error("Sending message failed", "extension", p.Name(), "error", err)
			}

		case "pgvissingup":
			if ok := tgs.CheckTargets(message, p.data.Targets); !ok {
				return errors.New("invalid target")
			}

			msgConfig := tgbotapi.NewMessage(
				message.Chat.ID,
				fmt.Sprintf(
					"Wenn du auf den \"Sing Up\" Button klickst, "+
						"bekommst du deinen Api Key zugesendet.\n\n"+
						"Bitte ignorieren, immer noch am testen.",
				),
			)

			msgConfig.ReplyToMessageID = message.MessageID

			button := tgbotapi.NewInlineKeyboardButtonURL("Sing Up", "t.me/talice_bot?start=pgvissingup")

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
