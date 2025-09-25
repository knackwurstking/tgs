package journal

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/templates"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Unit struct {
	// Name or the unit file to get
	Name string `json:"name" yaml:"name"`
	// Output will be used for the shell command `journalctl` as `--output ${output}`
	//
	// optional
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
}

type Units struct {
	System []Unit `json:"system,omitempty" yaml:"system,omitempty"`
	User   []Unit `json:"user,omitempty" yaml:"user,omitempty"`
}

func NewUnits() *Units {
	return &Units{
		System: make([]Unit, 0),
		User:   make([]Unit, 0),
	}
}

func (u *Units) GetSystemUnit(name string) (*Unit, error) {
	for i, unit := range u.System {
		if unit.Name == name {
			return &u.System[i], nil
		}
	}

	return nil, fmt.Errorf("user unit %s not found", name)
}

func (u *Units) GetUserUnit(name string) (*Unit, error) {
	for i, unit := range u.User {
		if unit.Name == name {
			return &u.User[i], nil
		}
	}

	return nil, fmt.Errorf("system unit %s not found", name)
}

func (u *Units) GetOutput(name string) (data []byte, err error) {
	isUser := true

	unit, err := u.GetUserUnit(name)
	if err != nil {
		isUser = false

		unit, err = u.GetSystemUnit(name)
		if err != nil {
			return nil, err
		}
	}

	var output string
	if unit.Output == "default" || unit.Output == "" {
		output = "short"
	} else {
		output = unit.Output
	}

	var cmd *exec.Cmd
	if isUser {
		cmd = exec.Command("bash", "-c",
			fmt.Sprintf(
				"journalctl --user -u %s --output %s --reverse --no-pager | sed 's/\x1b\\[[0-9;]*m//g'",
				unit.Name, output,
			),
		)
	} else {
		cmd = exec.Command("journalctl",
			"-u", unit.Name,
			"--output", output,
			"--reverse",
			"--no-pager",
		)
	}

	log.Debug("Executing journalctl command",
		"unit_name", unit.Name,
		"output_format", output,
		"is_user_service", isUser,
		"command", strings.Join(cmd.Args, " "),
	)

	if data, err = cmd.CombinedOutput(); err != nil {
		log.Error("Journalctl command failed",
			"unit_name", unit.Name,
			"command", strings.Join(cmd.Args, " "),
			"error", err,
		)
		return nil, err
	}

	log.Debug("Journalctl command executed successfully",
		"unit_name", unit.Name,
		"output_size_bytes", len(data),
	)
	return data, nil
}

type TemplateData struct {
	PageTitle   string
	SystemUnits []Unit
	UserUnits   []Unit
}

func (td *TemplateData) Patterns() []string {
	return []string{
		"data/index.go.html",
		"data/journallist.go.html", // block: content
		"data/ui-v4.3.0.css",       // block: ui
		"data/styles.css",          // block: styles
	}
}

type Data struct {
	Targets *tgs.Targets `yaml:"targets,omitempty"`
	Scopes  []tgs.Scope  `yaml:"scopes,omitempty"`
	Units   *Units       `yaml:"units,omitempty"`
}

type Journal struct {
	*tgbotapi.BotAPI

	data      *Data
	callbacks tgs.ReplyCallbacks
}

func New(api *tgbotapi.BotAPI) *Journal {
	return &Journal{
		BotAPI: api,
		data: &Data{
			Targets: tgs.NewTargets(),
			Scopes:  make([]tgs.Scope, 0),
			Units:   NewUnits(),
		},
		callbacks: tgs.ReplyCallbacks{},
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (j *Journal) Name() string {
	return "Journal"
}

func (j *Journal) SetBot(api *tgbotapi.BotAPI) {
	j.BotAPI = api
}

func (j *Journal) ConfigPath() string {
	return "journal.yaml"
}

func (j *Journal) MarshalYAML() (any, error) {
	return j.data, nil
}

func (j *Journal) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(j.data)
}

func (j *Journal) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/journal", "Get a journalctl log", j.data.Scopes)
	mbc.Add("/journallist", "List journalctl logs", j.data.Scopes)
}

func (j *Journal) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	if update.Message.ReplyToMessage != nil {
		if _, ok := j.callbacks.Get(update.Message.ReplyToMessage.MessageID); ok {
			return true
		}
	}

	return strings.HasPrefix(update.Message.Command(), "journal")
}

func (j *Journal) Handle(update tgbotapi.Update) error {
	if j.BotAPI == nil {
		log.Error("Journal extension BotAPI is nil")
		panic("BotAPI is nil!")
	}

	message := update.Message
	log.Debug("Journal extension handling update",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
		"command", message.Command(),
	)

	if ok := tgs.CheckTargets(message, j.data.Targets); !ok {
		log.Debug("Journal request from unauthorized target",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
			"command", message.Command(),
		)
		return errors.New("invalid target")
	}

	if message.ReplyToMessage != nil {
		replyMessageID := message.ReplyToMessage.MessageID
		if cb, ok := j.callbacks.Get(replyMessageID); ok {
			log.Debug("Processing journal reply callback",
				"reply_to_message_id", replyMessageID,
				"user_text", message.Text,
			)

			err := cb(message)
			if err != nil {
				log.Error("Journal reply callback failed",
					"reply_to_message_id", replyMessageID,
					"error", err,
				)

				msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("`error: %s`", err))
				msg.ParseMode = "MarkdownV2"
				msg.ReplyToMessageID = replyMessageID
				_, err = j.Send(msg)
			}

			return err
		}
	}

	switch command := message.Command(); command {
	case "journallist":
		log.Info("Generating journal units list",
			"user_id", message.From.ID,
			"system_units_count", len(j.data.Units.System),
			"user_units_count", len(j.data.Units.User),
		)

		content, err := templates.GetTemplateData(&TemplateData{
			PageTitle:   "Journal Units",
			SystemUnits: j.data.Units.System,
			UserUnits:   j.data.Units.User,
		})
		if err != nil {
			log.Error("Failed to generate journal units template",
				"error", err,
			)
			return err
		}

		documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
			Name:  "journal-units.html",
			Bytes: content,
		})
		documentConfig.ReplyToMessageID = message.MessageID

		log.Debug("Sending journal units list document",
			"document_size_bytes", len(content),
		)

		_, err = j.Send(documentConfig)
		if err != nil {
			log.Error("Failed to send journal units list",
				"error", err,
			)
		}
		return err
	case "journal":
		log.Info("Initiating journal request",
			"user_id", message.From.ID,
			"username", message.From.UserName,
		)

		msgConfig := tgbotapi.NewMessage(
			message.Chat.ID,
			"Hey there\\! Can you send me the name of the journal?\n\n"+
				"You'll have about 5 minutes to respond to this message\\.\n\n"+
				">You need to reply to this message for this to work\\.",
		)
		msgConfig.ReplyToMessageID = message.MessageID
		msgConfig.ParseMode = "MarkdownV2"

		msg, err := j.Send(msgConfig)
		if err != nil {
			log.Error("Failed to send journal request prompt",
				"error", err,
			)
			return err
		}

		j.callbacks.Set(msg.MessageID, j.replyCallbackJournalCommand)

		log.Debug("Journal callback registered",
			"callback_message_id", msg.MessageID,
			"timeout_minutes", 5,
		)

		go func() { // Auto Delete Function
			time.Sleep(time.Minute * 5)
			j.callbacks.Delete(msg.MessageID)
			log.Debug("Journal callback expired and removed",
				"callback_message_id", msg.MessageID,
			)
		}()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

func (j *Journal) replyCallbackJournalCommand(message *tgbotapi.Message) error {
	log.Info("Processing journal unit request",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"requested_text", message.Text,
		"reply_to_message_id", message.ReplyToMessage.MessageID,
		"message_id", message.MessageID,
	)

	messageSplit := strings.Split(message.Text, " ")
	for x := range messageSplit {
		messageSplit[x] = strings.ToLower(
			strings.Trim(messageSplit[x], " \r\n\t"),
		)
	}

	level := ""
	for _, t := range messageSplit {
		switch t {
		case "system":
			level = "system"
		case "user":
			level = "user"
		}
	}

	var (
		fileName string
		content  = []byte{}
		err      error
	)

	if level == "user" || level == "" {
		log.Debug("Searching user units",
			"search_terms", messageSplit,
			"available_units", len(j.data.Units.User),
		)

		for _, unit := range j.data.Units.User {
			if slices.Contains(messageSplit, strings.ToLower(unit.Name)) {
				log.Info("Found matching user unit",
					"unit_name", unit.Name,
					"unit_output", unit.Output,
				)
				content, err = j.data.Units.GetOutput(unit.Name)
				fileName = fmt.Sprintf("%s.log", unit.Name)
				break
			}
		}
	}

	if level == "system" || level == "" {
		log.Debug("Searching system units",
			"search_terms", messageSplit,
			"available_units", len(j.data.Units.System),
		)

		for _, unit := range j.data.Units.System {
			if slices.Contains(messageSplit, strings.ToLower(unit.Name)) {
				log.Info("Found matching system unit",
					"unit_name", unit.Name,
					"unit_output", unit.Output,
				)
				content, err = j.data.Units.GetOutput(unit.Name)
				fileName = fmt.Sprintf("%s.log", unit.Name)
				break
			}
		}
	}

	if err != nil {
		log.Error("Failed to get journal output",
			"error", err,
			"search_terms", strings.Join(messageSplit, " "),
		)
		return err
	}

	if len(content) == 0 {
		log.Warn("No journal unit found matching request",
			"search_terms", strings.Join(messageSplit, " "),
			"available_user_units", len(j.data.Units.User),
			"available_system_units", len(j.data.Units.System),
		)
		return fmt.Errorf("unit not found: %s", strings.Join(messageSplit, " "))
	}

	log.Info("Sending journal logs",
		"unit_file_name", fileName,
		"log_size_bytes", len(content),
		"user_id", message.From.ID,
	)

	documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: content,
	})
	documentConfig.ReplyToMessageID = message.MessageID

	_, err = j.Send(documentConfig)
	if err != nil {
		log.Error("Failed to send journal document",
			"file_name", fileName,
			"error", err,
		)
	}
	return err
}
