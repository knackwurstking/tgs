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
	mbc.Add("/pgpressregister", "Get an api key for the \"PG Presse Server\" project.", p.data.Scopes)
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
		log.Error("PGPress extension BotAPI is nil")
		panic("BotAPI is nil!")
	}

	message := update.Message
	log.Debug("PGPress extension handling update",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
		"command", message.Command(),
		"message_text", message.Text,
	)

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
	log.Info("Processing PGPress registration start",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"message_text", message.Text,
	)

	key := strings.SplitN(message.Text, "-", 2)[1]
	log.Debug("Extracted registration key from deep link", "key", key)

	if !slices.Contains(p.keys, key) {
		log.Warn("Invalid or expired registration key",
			"key", key,
			"user_id", message.From.ID,
			"available_keys_count", len(p.keys),
		)
		msgConfig := tgbotapi.NewMessage(message.From.ID,
			"Tut mir leid, Aber dieser \"Deep Link\" ist abgelaufen!")

		if _, err := p.Send(msgConfig); err != nil {
			log.Error("Failed to send expired link message",
				"user_id", message.From.ID,
				"error", err,
			)
		}

		return errors.New("invalid target")
	}

	userName := message.From.UserName
	if userName == "" {
		userName = strings.Trim(message.From.FirstName+" "+message.From.LastName, " ")
		log.Debug("Using fallback username",
			"user_id", message.From.ID,
			"fallback_name", userName,
		)
	}

	log.Debug("Creating PGPress user account",
		"user_id", message.From.ID,
		"username", userName,
	)

	user, err := NewUser(message.From.ID, userName)
	if err != nil {
		log.Error("Failed to create PGPress user",
			"user_id", message.From.ID,
			"username", userName,
			"error", err,
		)

		p.Send(tgbotapi.NewMessage(
			message.From.ID,
			fmt.Sprintf("Error: %s", err.Error()),
		))

		return err
	}

	log.Info("PGPress user created successfully",
		"user_id", message.From.ID,
		"username", userName,
		"api_key_length", len(user.ApiKey),
	)

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
		log.Error("Failed to send info message",
			"user_id", message.From.ID,
			"error", err,
		)
	}

	// ApiKey message
	msgConfig = tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("`%s`", user.ApiKey))
	msgConfig.ParseMode = "MarkdownV2"

	if _, err := p.Send(msgConfig); err != nil {
		log.Error("Failed to send API key message",
			"user_id", message.From.ID,
			"error", err,
		)
	} else {
		log.Info("API key sent to user",
			"user_id", message.From.ID,
		)
	}

	// Link to the pg-press server login page
	msgConfig = tgbotapi.NewMessage(message.From.ID,
		"Einfach den Api Key beim Login einfügen.")

	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonURL(
				"PG: Presse Server",
				"https://knackwurstking.com/pg-press/",
			),
		},
	)

	if _, err := p.Send(msgConfig); err != nil {
		log.Error("Failed to send login link message",
			"user_id", message.From.ID,
			"error", err,
		)
	} else {
		log.Debug("Login link sent to user",
			"user_id", message.From.ID,
		)
	}

	return nil
}

func (p *PGPress) handlePGPressRegister(message *tgbotapi.Message) error {
	log.Info("Processing PGPress registration command",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
	)

	if ok := tgs.CheckTargets(message, p.data.Targets); !ok {
		log.Debug("PGPress registration request from unauthorized target",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
		)
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
	log.Debug("Generated registration deep link",
		"key", key,
		"user_id", message.From.ID,
	)

	deepLink := fmt.Sprintf("t.me/talice_bot?start=pgpressregister-%s", key)
	button := tgbotapi.NewInlineKeyboardButtonURL("Registrieren", deepLink)
	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			button,
		},
	)

	if _, err := p.Send(msgConfig); err != nil {
		log.Error("Failed to send registration button",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
			"error", err,
		)
	} else {
		log.Debug("Registration button sent successfully",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
		)
	}

	p.keys = append(p.keys, key)
	log.Debug("Registration key stored",
		"key", key,
		"total_active_keys", len(p.keys),
	)

	return nil
}
