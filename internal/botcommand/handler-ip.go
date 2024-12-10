package botcommand

// TODO: Move to ip package, just like opmanga

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

const (
	IPURL = "https://ifconfig.io"
)

type IPConfig struct {
	Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
}

// IP implements the Handler interface
type IP struct {
	*tgbotapi.BotAPI
	targets  *Targets
	register []tgs.BotCommandScope
}

func NewIP(botAPI *tgbotapi.BotAPI) *IP {
	return &IP{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  NewTargets(),
	}
}

func (this *IP) MarshalJSON() ([]byte, error) {
	return json.Marshal(IPConfig{Register: this.register, Targets: this.targets})
}

func (this *IP) UnmarshalJSON(data []byte) error {
	d := IPConfig{Register: this.register, Targets: this.targets}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *IP) MarshalYAML() (interface{}, error) {
	return IPConfig{Register: this.register, Targets: this.targets}, nil
}

func (this *IP) UnmarshalYAML(value *yaml.Node) error {
	d := IPConfig{Register: this.register, Targets: this.targets}

	if err := value.Decode(&d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets

	return nil
}

func (this *IP) BotCommand() string {
	return "ip"
}

func (this *IP) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *IP) Targets() *Targets {
	return this.targets
}

func (this *IP) AddCommands(c *tgs.MyBotCommands) {
	c.Add("/"+this.BotCommand(), "Get server IP", this.Register())
}

func (this *IP) Run(message *tgbotapi.Message) error {
	address, err := this.fetchAddressFromURL()
	if err != nil {
		return err
	}

	msgConfig := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("`%s`", address))
	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = this.Send(msgConfig)
	return err
}

func (this *IP) fetchAddressFromURL() (address string, err error) {
	resp, err := http.Get(IPURL)
	if err != nil {
		return address, err
	}
	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("request to %s: %d (%s)",
			IPURL, resp.StatusCode, resp.Status,
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return address, err
	}

	return strings.Trim(string(data), "\n\r\t "), nil
}
