package botcommand

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type IP struct {
	*tgbotapi.BotAPI

	register []tgs.BotCommandScope
	targets  *Targets
}

func NewIP(botAPI *tgbotapi.BotAPI) *IP {
	return &IP{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  NewTargets(),
	}
}

func (this *IP) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Register []tgs.BotCommandScope `json:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	})
}

func (this *IP) MarshalYAML() (interface{}, error) {
	return struct {
		Register []tgs.BotCommandScope `yaml:"register,omitempty"`
		Targets  *Targets              `yaml:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	}, nil
}

func (this *IP) UnmarshalJSON(data []byte) error {
	d := struct {
		Register []tgs.BotCommandScope `json:"register,omitempty"`
		Targets  *Targets              `json:"targets,omitempty"`
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

func (this *IP) UnmarshalYAML(value *yaml.Node) error {
	d := struct {
		Register []tgs.BotCommandScope `yaml:"register,omitempty"`
		Targets  *Targets              `yaml:"targets,omitempty"`
	}{
		Register: this.register,
		Targets:  this.targets,
	}

	err := value.Decode(&d)
	if err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *IP) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *IP) Targets() *Targets {
	return this.targets
}

func (*IP) URL() string {
	return "https://ifconfig.io"
}

func (this *IP) Run(message *tgbotapi.Message) error {
	address, err := this.FetchAddressFromURL()
	if err != nil {
		return err
	}

	msgConfig := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("`%s`", address))
	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = this.BotAPI.Send(msgConfig)
	return err
}

func (this *IP) AddCommands(c *tgs.MyBotCommands, scopes ...tgs.BotCommandScope) {
	c.Add(BotCommandIP, "Get server IP", scopes)
}

func (this *IP) FetchAddressFromURL() (address string, err error) {
	resp, err := http.Get(this.URL())
	if err != nil {
		return address, err
	}
	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("request to %s: %d (%s)",
			this.URL(), resp.StatusCode, resp.Status,
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return address, err
	}

	return strings.Trim(string(data), "\n\r\t "), nil
}
