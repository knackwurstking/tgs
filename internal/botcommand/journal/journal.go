package journal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/templates"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

// Journal implements the Handler interface
type Journal struct {
	*tgbotapi.BotAPI
	targets  *botcommand.Targets
	units    *Units
	reply    chan *botcommand.Reply
	register []tgs.BotCommandScope
}

func NewJournal(botAPI *tgbotapi.BotAPI, reply chan *botcommand.Reply) *Journal {
	return &Journal{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),
		units:    NewUnits(),

		reply: reply,
	}
}

func (j *Journal) MarshalJSON() ([]byte, error) {
	return json.Marshal(JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	})
}

func (j *Journal) UnmarshalJSON(data []byte) error {
	d := JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	j.register = d.Register
	j.targets = d.Targets
	j.units = d.Units

	return nil
}

func (j *Journal) MarshalYAML() (interface{}, error) {
	return JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	}, nil
}

func (j *Journal) UnmarshalYAML(value *yaml.Node) error {
	d := JournalConfig{
		Register: j.register,
		Targets:  j.targets,
		Units:    j.units,
	}

	if err := value.Decode(&d); err != nil {
		return err
	}

	j.register = d.Register
	j.targets = d.Targets
	j.units = d.Units

	return nil
}

func (j *Journal) BotCommand() string {
	return "journal"
}

func (j *Journal) Register() []tgs.BotCommandScope {
	return j.register
}

func (j *Journal) Targets() *botcommand.Targets {
	return j.targets
}

func (j *Journal) AddCommands(c *tgs.MyBotCommands) {
	c.Add("/"+j.BotCommand()+"list", "List journalctl logs", j.Register())
	c.Add("/"+j.BotCommand(), "Get a journalctl log", j.Register())
}

func (j *Journal) Run(message *tgbotapi.Message) error {
	if j.isListCommand(message.Command()) {
		return j.handleListCommand(message)
	}

	msgConfig := tgbotapi.NewMessage(
		message.Chat.ID,
		"Hi there! Could you please send me the name of the journal? You’ll need to reply to this message to get it. I’ll give you a little time, about 5 minutes, to answer!",
	)
	msgConfig.ReplyToMessageID = message.MessageID

	msg, err := j.Send(msgConfig)
	if err != nil || j.reply == nil {
		return err
	}

	j.reply <- &botcommand.Reply{
		Message:  &msg,
		Timeout:  time.Minute * 5,
		Callback: j.replyCallback,
	}

	return nil
}

func (j *Journal) isListCommand(command string) bool {
	return command == j.BotCommand()+"list"
}

func (j *Journal) handleListCommand(message *tgbotapi.Message) error {
	content, err := templates.GetTemplateData(&TemplateData{
		PageTitle:   "Journal Units",
		SystemUnits: j.units.System,
		UserUnits:   j.units.User,
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
}

func (j *Journal) replyCallback(message *tgbotapi.Message) error {
	slog.Debug("Handle reply callback",
		"command", j.BotCommand(),
		"message.MessageID", message.MessageID,
		"message.Text", message.Text,
	)

	messageSplit := strings.Split(message.Text, " ")
	for x := 0; x < len(messageSplit); x++ {
		messageSplit[x] = strings.ToLower(
			strings.Trim(messageSplit[x], " \r\n\t"),
		)
	}

	level := "user"
	for _, t := range messageSplit {
		if t == "system" {
			level = "system"
		}
	}

	var (
		fileName string
		content  []byte
		err      error
	)

	if level == "system" {
		for _, unit := range j.units.System {
			if slices.Contains(messageSplit, strings.ToLower(unit.Name)) {
				content, err = j.units.GetOutput(unit.Name)
				fileName = fmt.Sprintf("%s.log", unit.Name)
				break
			}
		}
	}

	if level == "user" {
		for _, unit := range j.units.User {
			if slices.Contains(messageSplit, strings.ToLower(unit.Name)) {
				content, err = j.units.GetOutput(unit.Name)
				fileName = fmt.Sprintf("%s.log", unit.Name)
				break
			}
		}
	}

	if err != nil {
		return err
	}

	documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: content,
	})
	documentConfig.ReplyToMessageID = message.MessageID

	_, err = j.Send(documentConfig)
	return err
}
