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

	log.Debugf("Journal: Units: Run journalctl command with args %#v", cmd.Args)
	if data, err = cmd.CombinedOutput(); err != nil {
		return nil, err
	} else {
		return data, nil
	}
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
		panic("BotAPI is nil!")
	}

	message := update.Message

	if ok := tgs.CheckTargets(message, j.data.Targets); !ok {
		return errors.New("invalid target")
	}

	if message.ReplyToMessage != nil {
		replyMessageID := message.ReplyToMessage.MessageID
		if cb, ok := j.callbacks.Get(replyMessageID); ok {
			err := cb(message)
			if err != nil {
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
		content, err := templates.GetTemplateData(&TemplateData{
			PageTitle:   "Journal Units",
			SystemUnits: j.data.Units.System,
			UserUnits:   j.data.Units.User,
		})
		if err != nil {
			return err
		}

		documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
			Name:  "journal-units.html",
			Bytes: content,
		})
		documentConfig.ReplyToMessageID = message.MessageID

		_, err = j.Send(documentConfig)
		return err
	case "journal":
		msgConfig := tgbotapi.NewMessage(
			message.Chat.ID,
			"Hey there\\! Can you send me the name of the journal?\n\n"+
				"You’ll have about 5 minutes to respond to this message\\.\n\n"+
				">You need to reply to this message for this to work\\.",
		)
		msgConfig.ReplyToMessageID = message.MessageID
		msgConfig.ParseMode = "MarkdownV2"

		msg, err := j.Send(msgConfig)
		if err != nil {
			return err
		}

		j.callbacks.Set(msg.MessageID, j.replyCallbackJournalCommand)

		go func() { // Auto Delete Function
			time.Sleep(time.Minute * 5)
			j.callbacks.Delete(msg.MessageID)
		}()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

func (j *Journal) replyCallbackJournalCommand(message *tgbotapi.Message) error {
	log.Debugf(
		"Journal: Handle reply callback for \"%s\" (%d): text=%s; message_id=%d",
		message.Command(), message.ReplyToMessage.MessageID,
		message.Text,
		message.MessageID,
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
		for _, unit := range j.data.Units.User {
			if slices.Contains(messageSplit, strings.ToLower(unit.Name)) {
				content, err = j.data.Units.GetOutput(unit.Name)
				fileName = fmt.Sprintf("%s.log", unit.Name)
				break
			}
		}
	}

	if level == "system" || level == "" {
		for _, unit := range j.data.Units.System {
			if slices.Contains(messageSplit, strings.ToLower(unit.Name)) {
				content, err = j.data.Units.GetOutput(unit.Name)
				fileName = fmt.Sprintf("%s.log", unit.Name)
				break
			}
		}
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return fmt.Errorf("unit not found: %s", strings.Join(messageSplit, " "))
	}

	documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: content,
	})
	documentConfig.ReplyToMessageID = message.MessageID

	_, err = j.Send(documentConfig)
	return err
}
