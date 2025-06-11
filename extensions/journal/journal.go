package ip

import (
	"fmt"
	"log/slog"
	"os/exec"
	"slices"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/templates"
	"github.com/knackwurstking/tgs/pkg/extension"
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

	slog.Debug("Run journalctl command", "args", cmd.Args)
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
		"data/ui-v2.0.0.css",       // block: ui
		"data/styles.css",          // block: styles
	}
}

type Data struct {
	Targets  *botcommand.Targets    `yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope  `yaml:"register,omitempty"`
	Units    *Units                 `yaml:"units,omitempty"`
	Reply    chan *botcommand.Reply `yaml:"-"`
}

type Journal struct {
	*tgbotapi.BotAPI

	data *Data
}

func New(api *tgbotapi.BotAPI) *Journal {
	return &Journal{
		BotAPI: api,
		data: &Data{
			Targets:  botcommand.NewTargets(),
			Register: make([]tgs.BotCommandScope, 0),
			Units:    NewUnits(),
			// Reply: ,
		},
	}
}

func NewExtension(api *tgbotapi.BotAPI) extension.Extension {
	return New(api)
}

func (j *Journal) SetBot(api *tgbotapi.BotAPI) {
	j.BotAPI = api
}

func (j *Journal) ConfigPath() string {
	return "journal.config"
}

func (j *Journal) MarshalYAML() (any, error) {
	return j.data, nil
}

func (j *Journal) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(j.data)
}

func (j *Journal) Register() []tgs.BotCommandScope {
	return j.data.Register
}

func (j *Journal) Targets() *botcommand.Targets {
	return j.data.Targets
}

func (j *Journal) Commands(mbc *tgs.MyBotCommands) {
	mbc.Add("/journal", "Get a journalctl log", j.Register())
	mbc.Add("/journallist", "List journalctl logs", j.Register())
}

func (j *Journal) Is(command string) bool {
	return strings.HasPrefix(command, "journal")
}

func (j *Journal) Handle(message *tgbotapi.Message) error {
	if j.BotAPI == nil {
		panic("BotAPI is nil!")
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
		if err != nil || j.data.Reply == nil {
			return err
		}

		// FIXME: Reply is not nil for now, need to implement this somehow
		j.data.Reply <- &botcommand.Reply{
			Message:  &msg,
			Timeout:  time.Minute * 5,
			Callback: j.replyCallback,
		}
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

func (j *Journal) replyCallback(message *tgbotapi.Message) error {
	slog.Debug("Handle reply callback",
		"command", message.Command(),
		"message.MessageID", message.MessageID,
		"message.Text", message.Text,
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
