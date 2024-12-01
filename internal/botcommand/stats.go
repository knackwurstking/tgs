package botcommand

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type Stats struct {
	*tgbotapi.BotAPI

	register []tgs.BotCommandScope
	targets  *Targets
}

func NewStats(botAPI *tgbotapi.BotAPI) *Stats {
	return &Stats{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  NewTargets(),
	}
}

func (this *Stats) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	})
}

func (this *Stats) MarshalYAML() (interface{}, error) {
	return struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	}, nil
}

func (this *Stats) UnmarshalJSON(data []byte) error {
	d := struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	}

	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *Stats) UnmarshalYAML(value *yaml.Node) error {
	d := struct {
		Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	}

	err := value.Encode(&d)
	if err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *Stats) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *Stats) Targets() *Targets {
	return this.targets
}

func (this *Stats) Run(message *tgbotapi.Message) error {
	data := struct {
		UserName        string `json:"username"`
		UserID          int64  `json:"user_id"`
		ChatID          int64  `json:"chat_id"`
		MessageThreadID int    `json:"message_thread_id"`
	}{
		UserName:        message.From.UserName,
		UserID:          message.From.ID,
		ChatID:          message.Chat.ID,
		MessageThreadID: message.MessageThreadID,
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")

	msgConfig := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf("```json\n")+
			fmt.Sprintf("%s\n", string(jsonData))+
			fmt.Sprintf("```"),
	)

	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = this.BotAPI.Send(msgConfig)
	return err
}

func (this *Stats) AddCommands(c *tgs.MyBotCommands, scopes ...tgs.BotCommandScope) {
	c.Add(BotCommandStats, "Get ID info", scopes)
}
