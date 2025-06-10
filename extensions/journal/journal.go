package ip

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/extension"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Data struct {
	Targets  *botcommand.Targets   `yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `yaml:"register,omitempty"`
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
		},
	}
}

func NewExtension(api *tgbotapi.BotAPI) extension.Extension {
	return New(api)
}

func (j *Journal) SetBot(api *tgbotapi.BotAPI) {
	j.BotAPI = api
}

// TODO: Continue here
func (ip *IP) ConfigPath() string {
	return "ip.config"
}

func (ip *IP) MarshalYAML() (any, error) {
	return ip.data, nil
}

func (ip *IP) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(ip.data); err != nil {
		return err
	}

	return nil
}

func (ip *IP) Register() []tgs.BotCommandScope {
	return ip.data.Register
}

func (ip *IP) Targets() *botcommand.Targets {
	return ip.data.Targets
}

func (ip *IP) Commands(mbc *tgs.MyBotCommands) {
	mbc.Add("/ip", "Get server IP", ip.Register())
}

func (ip *IP) Is(command string) bool {
	return strings.HasPrefix(command, "ip")
}

func (ip *IP) Handle(message *tgbotapi.Message) error {
	if ip.BotAPI == nil {
		panic("BotAPI is nil!")
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
