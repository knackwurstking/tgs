package stats

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

// Stats implements the Handler interface
type Stats struct {
	*tgbotapi.BotAPI
	targets  *botcommand.Targets
	register []tgs.BotCommandScope
}

func NewStats(botAPI *tgbotapi.BotAPI) *Stats {
	return &Stats{
		BotAPI: botAPI,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),
	}
}

func (s *Stats) MarshalJSON() ([]byte, error) {
	return json.Marshal(StatsConfig{Register: s.register, Targets: s.targets})
}

func (s *Stats) UnmarshalJSON(data []byte) error {
	d := StatsConfig{Register: s.register, Targets: s.targets}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	s.register = d.Register
	s.targets = d.Targets

	return nil
}

func (s *Stats) MarshalYAML() (interface{}, error) {
	return StatsConfig{Register: s.register, Targets: s.targets}, nil
}

func (s *Stats) UnmarshalYAML(value *yaml.Node) error {
	d := StatsConfig{Register: s.register, Targets: s.targets}

	if err := value.Decode(&d); err != nil {
		return err
	}

	s.register = d.Register
	s.targets = d.Targets

	return nil
}

func (s *Stats) BotCommand() string {
	return "stats"
}

func (s *Stats) Register() []tgs.BotCommandScope {
	return s.register
}

func (s *Stats) Targets() *botcommand.Targets {
	return s.targets
}

func (s *Stats) AddCommands(c *tgs.MyBotCommands) {
	c.Add("/"+s.BotCommand(), "Get ID info", s.Register())
}

func (s *Stats) Run(message *tgbotapi.Message) error {
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
		"```json\n"+fmt.Sprintf("%s\n", string(jsonData))+"```",
	)

	msgConfig.ReplyToMessageID = message.MessageID
	msgConfig.ParseMode = "MarkdownV2"

	_, err = s.Send(msgConfig)
	return err
}
