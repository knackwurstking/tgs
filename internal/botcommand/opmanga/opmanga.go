package opmanga

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/knackwurstking/tgs/internal/botcommand"
	"github.com/knackwurstking/tgs/internal/templates"
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

// OPManga implements the Handler interface
type OPManga struct {
	*tgbotapi.BotAPI
	targets  *botcommand.Targets
	reply    chan *botcommand.Reply
	path     string
	register []tgs.BotCommandScope
}

func NewOPManga(bot *tgbotapi.BotAPI) *OPManga {
	return &OPManga{
		BotAPI: bot,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),
	}
}

func (this *OPManga) MarshalJSON() ([]byte, error) {
	return json.Marshal(Config{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	})
}

func (this *OPManga) UnmarshalJSON(data []byte) error {
	d := Config{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets
	this.path = d.Path

	return nil
}

func (this *OPManga) MarshalYAML() (interface{}, error) {
	return Config{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}, nil
}

func (this *OPManga) UnmarshalYAML(value *yaml.Node) error {
	d := Config{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}

	if err := value.Decode(&d); err != nil {
		return err
	}

	this.register = d.Register
	this.targets = d.Targets
	this.path = d.Path

	return nil
}

func (this *OPManga) BotCommand() string {
	return "opmanga"
}

func (this *OPManga) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *OPManga) Targets() *botcommand.Targets {
	return this.targets
}

func (this *OPManga) AddCommands(c *tgs.MyBotCommands) {
	c.Add("/"+this.BotCommand()+"list", "List all available chapters", this.Register())
	c.Add("/"+this.BotCommand(), "Request a chapter", this.Register())
}

func (this *OPManga) Run(m *tgbotapi.Message) error {
	if this.isListCommand(m.Command()) {
		return this.handleListCommand(m)
	}

	msgConfig := tgbotapi.NewMessage(
		m.Chat.ID,
		"Hi there! Reply to this message to get the episode you want. I’ll give you about 5 minutes to answer.",
	)
	msgConfig.ReplyToMessageID = m.MessageID

	msg, err := this.Send(msgConfig)
	if err != nil || this.reply == nil {
		return err
	}

	this.reply <- &botcommand.Reply{
		Message:  &msg,
		Timeout:  time.Minute * 5,
		Callback: this.replyCallback,
	}

	return nil
}

func (this *OPManga) arcs() ([]Arc, error) {
	if this.path == "" {
		return nil, fmt.Errorf("missing path")
	}

	info, err := os.Stat(this.path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("nope, need an directory here")
	}

	dirEntries, err := os.ReadDir(this.path)
	if err != nil {
		return nil, err
	}

	arcs := make([]Arc, 0)

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue // Just ignore all non directories
		}

		sub, err := os.ReadDir(filepath.Join(this.path, dirEntry.Name()))
		if err != nil {
			continue // Ignore for now
		}

		arc := Arc{
			Chapters: []*Chapter{},
		}

		if s := strings.SplitN(dirEntry.Name(), " ", 2); len(s) < 2 {
			arc.Name = dirEntry.Name()
		} else {
			arc.Name = s[1] // Ignore the prefixed number (ex.: "016 Thousand Sunny Arc")
		}

		for _, subEntry := range sub {
			if subEntry.IsDir() {
				continue // Skip all directories
			}

			if filepath.Ext(subEntry.Name()) != ".pdf" {
				continue // Allow only pdf
			}

			chapter, err := NewChapter(
				filepath.Join(this.path, dirEntry.Name(), subEntry.Name()),
			)
			if err != nil {
				return nil, err
			}

			arc.Chapters = append(arc.Chapters, chapter)
		}

		arcs = append(arcs, arc)
	}

	return arcs, nil
}

func (this *OPManga) isListCommand(c string) bool {
	return c == this.BotCommand()+"list"
}

func (this *OPManga) handleListCommand(m *tgbotapi.Message) error {
	arcs, err := this.arcs()
	if err != nil {
		return err
	}

	content, err := templates.GetTemplateData(&TemplateData{
		PageTitle: "One Piece Manga",
		Arcs:      arcs,
	})
	if err != nil {
		return err
	}

	documentConfig := tgbotapi.NewDocument(m.Chat.ID, tgbotapi.FileBytes{
		Name:  "opmanga-chapters.html",
		Bytes: content,
	})
	documentConfig.ReplyToMessageID = m.MessageID

	_, err = this.Send(documentConfig)
	return err
}

func (this *OPManga) replyCallback(message *tgbotapi.Message) error {
	slog.Debug("Handle reply callback",
		"command", this.BotCommand(),
		"message.MessageID", message.MessageID,
		"message.Text", message.Text,
	)

	arcs, err := this.arcs()
	if err != nil {
		return err
	}

	// Parse message and get episode string
	match := regexp.MustCompile(`[0-9]+`).FindString(message.Text)
	if match == "" {
		return fmt.Errorf("nothing found, need a number here: %s", message.Text)
	}

	n, err := strconv.Atoi(match)
	if err != nil {
		return err
	}

	// Search arcs data for chapter
outer_loop:
	for _, a := range arcs {
		for _, c := range a.Chapters {
			if c.Number() == n {
				if pdf, err := c.PDF(); err != nil {
					return err
				} else {
					chatID := message.Chat.ID

					documentConfig := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
						Name:  pdf.Name(),
						Bytes: pdf.Data(),
					})

					msgConfig := tgbotapi.NewMessage(chatID, "")
					msgConfig.ReplyToMessageID = message.ReplyToMessage.MessageID
					msgConfig.ReplyMarkup = documentConfig

					_, err = this.Send(msgConfig)
					if err != nil {
						return err
					}
				}

				break outer_loop
			}
		}
	}

	return fmt.Errorf("chapter number %d not found", n)
}
