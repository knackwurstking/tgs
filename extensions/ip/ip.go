package ip

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
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
		log.Error("IP extension BotAPI is nil")
		panic("BotAPI is nil!")
	}

	message := update.Message
	log.Debug("IP extension handling update",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
		"command", message.Command(),
	)

	if ok := tgs.CheckTargets(message, ip.data.Targets); !ok {
		log.Debug("IP request from unauthorized target",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
			"command", message.Command(),
		)
		return errors.New("invalid target")
	}

	if command := message.Command(); command != "ip" {
		log.Warn("Unknown command in IP extension",
			"command", command,
			"user_id", message.From.ID,
		)
		return fmt.Errorf("unknown command: %s", command)
	}

	log.Info("Processing IP lookup request",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
	)

	log.Debug("Starting IP address lookups")

	ipv4, err := ip.GetIPv4AddressFromURL()
	if err != nil {
		log.Warn("IPv4 lookup failed", "error", err)
		ipv4 = err.Error()
	} else {
		log.Debug("IPv4 lookup successful", "address", ipv4)
	}

	ipv6, err := ip.GetIPv6AddressFromURL()
	if err != nil {
		log.Warn("IPv6 lookup failed", "error", err)
		ipv6 = err.Error()
	} else {
		log.Debug("IPv6 lookup successful", "address", ipv6)
	}

	msgConfig := tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf(
			"**IPv4**: `%s`\n**IPv6**: `%s`", ipv4, ipv6,
		),
	)
	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	log.Debug("Sending IP information response",
		"chat_id", message.Chat.ID,
		"reply_to_message_id", message.MessageID,
	)

	_, err = ip.Send(msgConfig)
	if err != nil {
		log.Error("Failed to send IP information",
			"chat_id", message.Chat.ID,
			"error", err,
		)
	} else {
		log.Info("IP information sent successfully",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
		)
	}
	return err
}

func (ip *IP) GetIPv4AddressFromURL() (address string, err error) {
	start := time.Now()
	log.Debug("Starting IPv4 lookup", "url", IP4URL)

	resp, err := http.Get(IP4URL)
	if err != nil {
		log.Error("IPv4 HTTP request failed",
			"url", IP4URL,
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return address, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("IPv4 lookup returned non-OK status",
			"url", IP4URL,
			"status_code", resp.StatusCode,
			"status", resp.Status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return address, fmt.Errorf("request to %s: %d (%s)",
			IP4URL, resp.StatusCode, resp.Status,
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read IPv4 response body",
			"url", IP4URL,
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return address, err
	}

	address = strings.Trim(string(data), "\n\r\t ")
	log.Debug("IPv4 lookup completed successfully",
		"url", IP4URL,
		"address", address,
		"response_size", len(data),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return address, nil
}

func (ip *IP) GetIPv6AddressFromURL() (address string, err error) {
	start := time.Now()
	log.Debug("Starting IPv6 lookup", "url", IP6URL)

	resp, err := http.Get(IP6URL)
	if err != nil {
		log.Error("IPv6 HTTP request failed",
			"url", IP6URL,
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return address, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("IPv6 lookup returned non-OK status",
			"url", IP6URL,
			"status_code", resp.StatusCode,
			"status", resp.Status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return address, fmt.Errorf("request to %s: %d (%s)",
			IP6URL, resp.StatusCode, resp.Status,
		)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read IPv6 response body",
			"url", IP6URL,
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return address, err
	}

	address = strings.Trim(string(data), "\n\r\t ")
	log.Debug("IPv6 lookup completed successfully",
		"url", IP6URL,
		"address", address,
		"response_size", len(data),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return address, nil
}
