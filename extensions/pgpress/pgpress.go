package pgpress

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Data struct {
	Targets *tgs.Targets `yaml:"targets,omitempty"`
	Scopes  []tgs.Scope  `yaml:"scopes,omitempty"`
}

type PGPress struct {
	*tgbotapi.BotAPI

	data *Data

	keys []string
}

func New(api *tgbotapi.BotAPI) *PGPress {
	return &PGPress{
		data: &Data{
			Targets: tgs.NewTargets(),
			Scopes:  make([]tgs.Scope, 0),
		},
		keys: make([]string, 0),
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (p *PGPress) Name() string {
	return "pgpress"
}

func (p *PGPress) SetBot(api *tgbotapi.BotAPI) {
	p.BotAPI = api
}

func (p *PGPress) ConfigPath() string {
	return "pgpress.yaml"
}

func (p *PGPress) MarshalYAML() (any, error) {
	return p.data, nil
}

func (p *PGPress) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(p.data)
}

func (p *PGPress) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/pgpressregister", "Get an api key for the \"PG Vis Server\" project.", p.data.Scopes)
}

func (p *PGPress) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	command := update.Message.Command()

	if command == "start" {
		return strings.HasPrefix(update.Message.Text, "/start pgpressregister-")
	}

	return strings.HasPrefix(command, "pgpress")
}

func (p *PGPress) Handle(update tgbotapi.Update) error {
	if p.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	message := update.Message

	if message != nil {
		switch command := message.Command(); command {
		case "start":
			return p.handleStartPGPressRegister(message)

		case "pgpressregister":
			return p.handlePGPressRegister(message)

		default:
			return fmt.Errorf("unknown command: %s", command)
		}
	}

	return fmt.Errorf("there is nothing to do here")
}

func (p *PGPress) handleStartPGPressRegister(message *tgbotapi.Message) error {
	key := strings.SplitN(message.Text, "-", 2)[1]
	if !slices.Contains(p.keys, key) {
		msgConfig := tgbotapi.NewMessage(message.From.ID,
			"Tut mir leid, Aber dieser \"Deep Link\" ist abgelaufen!")

		if _, err := p.Send(msgConfig); err != nil {
			log.Errorf("PGPress: Send message failed: %s", err)
		}

		return errors.New("invalid target")
	}

	userName := message.From.UserName
	if userName == "" {
		userName = strings.Trim(message.From.FirstName+" "+message.From.LastName, " ")
	}

	user, err := NewUser(message.From.ID, userName)
	if err != nil {
		p.Send(tgbotapi.NewMessage(
			message.From.ID,
			fmt.Sprintf("Error: %s", err.Error()),
		))

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
		log.Errorf("PGPress: Send message failed: %s", err)
	}

	// ApiKey message
	msgConfig = tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("`%s`", user.ApiKey))
	msgConfig.ParseMode = "MarkdownV2"

	if _, err := p.Send(msgConfig); err != nil {
		log.Errorf("PGPress: Send message failed: %s", err)
	}

	// Link to the pg-vis server login page
	msgConfig = tgbotapi.NewMessage(message.From.ID,
		"Einfach den Api Key beim Login einfügen.")

	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonURL(
				"PG: Vis Server",
				"https://knackwurstking.com/pg-vis/",
			),
		},
	)

	if _, err := p.Send(msgConfig); err != nil {
		log.Errorf("PGPress: Send message failed: %s", err)
	}

	return nil
}

func (p *PGPress) handlePGPressRegister(message *tgbotapi.Message) error {
	if ok := tgs.CheckTargets(message, p.data.Targets); !ok {
		return errors.New("invalid target")
	}

	msgConfig := tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf(
			"Wenn du auf den \"Registrieren\" Button klickst, "+
				"bekommst du deinen Api Key zugesendet.",
		),
	)

	msgConfig.ReplyToMessageID = message.MessageID

	key := uuid.New().String()

	button := tgbotapi.NewInlineKeyboardButtonURL("Registrieren", fmt.Sprintf("t.me/talice_bot?start=pgpressregister-%s", key))
	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			button,
		},
	)

	if _, err := p.Send(msgConfig); err != nil {
		log.Errorf("PGPress: Send message failed: %s", err)
	}

	p.keys = append(p.keys, key)

	return nil
}
