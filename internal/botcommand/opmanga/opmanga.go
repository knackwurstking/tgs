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

func NewOPManga(bot *tgbotapi.BotAPI, reply chan *botcommand.Reply) *OPManga {
	return &OPManga{
		BotAPI: bot,

		register: []tgs.BotCommandScope{},
		targets:  botcommand.NewTargets(),

		reply: reply,
	}
}

func (opm *OPManga) MarshalJSON() ([]byte, error) {
	return json.Marshal(Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	})
}

func (opm *OPManga) UnmarshalJSON(data []byte) error {
	d := Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	opm.register = d.Register
	opm.targets = d.Targets
	opm.path = d.Path

	return nil
}

func (opm *OPManga) MarshalYAML() (interface{}, error) {
	return Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	}, nil
}

func (opm *OPManga) UnmarshalYAML(value *yaml.Node) error {
	d := Config{
		Register: opm.register,
		Targets:  opm.targets,
		Path:     opm.path,
	}

	if err := value.Decode(&d); err != nil {
		return err
	}

	opm.register = d.Register
	opm.targets = d.Targets
	opm.path = d.Path

	return nil
}

func (opm *OPManga) BotCommand() string {
	return "opmanga"
}

func (opm *OPManga) Register() []tgs.BotCommandScope {
	return opm.register
}

func (opm *OPManga) Targets() *botcommand.Targets {
	return opm.targets
}

func (opm *OPManga) AddCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/"+opm.BotCommand()+"list", "List all available chapters", opm.Register())
	mbc.Add("/"+opm.BotCommand(), "Request a chapter", opm.Register())
}

func (opm *OPManga) Run(m *tgbotapi.Message) error {
	if opm.isListCommand(m.Command()) {
		return opm.handleListCommand(m)
	}

	msgConfig := tgbotapi.NewMessage(
		m.Chat.ID,
		"Hey there! I need the chapter number. Reply to this message and I’ll send you the chapter. You’ll have about 5 minutes to respond.",
	)
	msgConfig.ReplyToMessageID = m.MessageID

	msg, err := opm.Send(msgConfig)
	if err != nil || opm.reply == nil {
		return err
	}

	opm.reply <- &botcommand.Reply{
		Message:  &msg,
		Timeout:  time.Minute * 5,
		Callback: opm.replyCallback,
	}

	return nil
}

func (opm *OPManga) arcs() ([]Arc, error) {
	if opm.path == "" {
		return nil, fmt.Errorf("missing path")
	}

	info, err := os.Stat(opm.path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("nope, need an directory here")
	}

	dirEntries, err := os.ReadDir(opm.path)
	if err != nil {
		return nil, err
	}

	arcs := make([]Arc, 0)

	for _, e := range dirEntries {
		if !e.IsDir() {
			continue // Just ignore all non directories
		}

		subDirEntries, err := os.ReadDir(filepath.Join(opm.path, e.Name()))
		if err != nil {
			continue // Ignore for now
		}

		arc := Arc{
			Chapters: []*Chapter{},
		}

		if s := strings.SplitN(e.Name(), " ", 2); len(s) < 2 {
			arc.Name = e.Name()
		} else {
			arc.Name = s[1] // Ignore the prefixed number (ex.: "016 Thousand Sunny Arc")
		}

		for _, e2 := range subDirEntries {
			if e2.IsDir() {
				continue // Skip all directories
			}

			if filepath.Ext(e2.Name()) != ".pdf" {
				continue // Allow only pdf
			}

			chapter, err := NewChapter(
				filepath.Join(opm.path, e.Name(), e2.Name()),
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

func (opm *OPManga) isListCommand(c string) bool {
	return c == opm.BotCommand()+"list"
}

func (opm *OPManga) handleListCommand(m *tgbotapi.Message) error {
	arcs, err := opm.arcs()
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

	_, err = opm.Send(documentConfig)
	return err
}

func (opm *OPManga) replyCallback(message *tgbotapi.Message) error {
	slog.Debug("Handle reply callback",
		"command", opm.BotCommand(),
		"message.MessageID", message.MessageID,
		"message.Text", message.Text,
	)

	arcs, err := opm.arcs()
	if err != nil {
		return err
	}

	// Parse message and get episode string
	match := regexp.MustCompile(`[0-9]+`).FindString(message.Text)
	if match == "" {
		return fmt.Errorf("nothing found, need a number here: %s", message.Text)
	}

	chapterNumber, err := strconv.Atoi(match)
	if err != nil {
		return err
	}

	// Search arcs data for chapter
	for _, a := range arcs {
		for _, c := range a.Chapters {
			if c.Number() == chapterNumber {
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

					_, err = opm.Send(msgConfig)
					if err != nil {
						slog.Error("Send chapter", "error", err)
						return err
					}

					return nil
				}
			}
		}
	}

	return fmt.Errorf("chapter number %d not found", chapterNumber)
}
