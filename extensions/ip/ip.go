package ip

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/pkg/tgs"
)

const (
	IP4URL = "https://ifconfig.io"
	IP6URL = "https://ipv6.icanhazip.com"
)

type Data struct {
	Targets *tgs.Targets `yaml:"targets,omitempty"`
	Scopes  []tgs.Scope  `yaml:"scopes,omitempty"`
}

type IP struct {
	*tgbotapi.BotAPI

	data *Data
}

func New(api *tgbotapi.BotAPI) *IP {
	return &IP{
		BotAPI: api,
		data: &Data{
			Targets: tgs.NewTargets(),
			Scopes:  make([]tgs.Scope, 0),
		},
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (ip *IP) Name() string {
	return "IP"
}

func (ip *IP) SetBot(api *tgbotapi.BotAPI) {
	ip.BotAPI = api
}

func (ip *IP) ConfigPath() string {
	return "ip.yaml"
}

func (ip *IP) MarshalYAML() (any, error) {
	return ip.data, nil
}

func (ip *IP) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(ip.data)
}

func (ip *IP) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/ip", "Get server IP", ip.data.Scopes)
}

func (ip *IP) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	return strings.HasPrefix(update.Message.Command(), "ip")
}

func (ip *IP) Handle(update tgbotapi.Update) error {
	if ip.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	message := update.Message

	if ok := tgs.CheckTargets(message, ip.data.Targets); !ok {
		return errors.New("invalid target")
	}

	if command := message.Command(); command != "ip" {
		return fmt.Errorf("unknown command: %s", command)
	}

	ipv4, err := ip.GetIPv4AddressFromURL()
	if err != nil {
		ipv4 = err.Error()
	}

	ipv6, err := ip.GetIPv6AddressFromURL()
	if err != nil {
		ipv6 = err.Error()
	}

	msgConfig := tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf(
			"**IPv4**: `%s`\n**IPv6**: `%s`", ipv4, ipv6,
		),
	)
	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = ip.Send(msgConfig)
	return err
}

func (ip *IP) GetIPv4AddressFromURL() (address string, err error) {
	resp, err := http.Get(IP4URL)
	if err != nil {
		return address, err
	}
	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("request to %s: %d (%s)",
			IP4URL, resp.StatusCode, resp.Status,
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return address, err
	}

	return strings.Trim(string(data), "\n\r\t "), nil
}

func (ip *IP) GetIPv6AddressFromURL() (address string, err error) {
	resp, err := http.Get(IP6URL)
	if err != nil {
		return address, err
	}
	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("request to %s: %d (%s)",
			IP6URL, resp.StatusCode, resp.Status,
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return address, err
	}

	return strings.Trim(string(data), "\n\r\t "), nil
}
