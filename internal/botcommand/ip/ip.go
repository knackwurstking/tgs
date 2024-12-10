package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

const (
	IPURL = "https://ifconfig.io"
)

// IP implements the Handler interface
type IP struct {
	*tgbotapi.BotAPI
	targets  *botcommand.Targets
	register []tgs.BotCommandScope
}

func NewIP(botAPI *tgbotapi.BotAPI) *IP {
	return &IP{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),
	}
}

func (ip *IP) MarshalJSON() ([]byte, error) {
	return json.Marshal(IPConfig{Register: ip.register, Targets: ip.targets})
}

func (ip *IP) UnmarshalJSON(data []byte) error {
	d := IPConfig{Register: ip.register, Targets: ip.targets}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	ip.register = d.Register
	ip.targets = d.Targets

	return nil
}

func (ip *IP) MarshalYAML() (interface{}, error) {
	return IPConfig{Register: ip.register, Targets: ip.targets}, nil
}

func (ip *IP) UnmarshalYAML(value *yaml.Node) error {
	d := IPConfig{Register: ip.register, Targets: ip.targets}

	if err := value.Decode(&d); err != nil {
		return err
	}

	ip.register = d.Register
	ip.targets = d.Targets

	return nil
}

func (ip *IP) BotCommand() string {
	return "ip"
}

func (ip *IP) Register() []tgs.BotCommandScope {
	return ip.register
}

func (ip *IP) Targets() *botcommand.Targets {
	return ip.targets
}

func (ip *IP) AddCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/"+ip.BotCommand(), "Get server IP", ip.Register())
}

func (ip *IP) Run(message *tgbotapi.Message) error {
	address, err := ip.fetchAddressFromURL()
	if err != nil {
		return err
	}

	msgConfig := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("`%s`", address))
	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = ip.Send(msgConfig)
	return err
}

func (ip *IP) fetchAddressFromURL() (address string, err error) {
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