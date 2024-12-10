package botcommand

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
	"github.com/knackwurstking/tgs/pkg/tgs"
	"gopkg.in/yaml.v3"
)

type PDF struct {
	name string
	data []byte
}

func NewPDF(name string, data []byte) *PDF {
	if data == nil {
		data = []byte{}
	}

	return &PDF{
		name: name,
		data: data,
	}
}

type OPMangaChapter struct {
	Name   string
	Path   string
	Number int
}

func (this *OPMangaChapter) PDF() (pdf File, err error) {
	// TODO: Read pdf data from path and return

	return nil, fmt.Errorf("under construction")
}

type OPMangaArc struct {
	Name     string
	Chapters []OPMangaChapter
}

type OPMangaTemplateData struct {
	PageTitle string
	Arcs      []OPMangaArc
}

func (this *OPMangaTemplateData) Patterns() []string {
	return []string{
		"templates/index.html",
		"templates/opmangalist.html",
		//"templates/pico.min.css",
		"templates/styles.css",
		"templates/original.css",
		"templates/ui.min.css",
		//"templates/ui.min.umd.cjs",
	}
}

type OPMangaConfig struct {
	Targets  *Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	Path     string                `json:"path" yaml:"path"`
	Register []tgs.BotCommandScope `json:"register,omitempty" yaml:"register,omitempty"`
}

// OPManga implements the Handler interface
type OPManga struct {
	*tgbotapi.BotAPI
	targets  *Targets
	reply    chan *Reply
	path     string
	register []tgs.BotCommandScope
}

func NewOPManga(bot *tgbotapi.BotAPI) *OPManga {
	return &OPManga{
		BotAPI: bot,

		register: []tgs.BotCommandScope{},
		targets:  NewTargets(),
	}
}

func (this *OPManga) MarshalJSON() ([]byte, error) {
	return json.Marshal(OPMangaConfig{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	})
}

func (this *OPManga) UnmarshalJSON(data []byte) error {
	d := OPMangaConfig{
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
	return OPMangaConfig{
		Register: this.register,
		Targets:  this.targets,
		Path:     this.path,
	}, nil
}

func (this *OPManga) UnmarshalYAML(value *yaml.Node) error {
	d := OPMangaConfig{
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

func (this *OPManga) Register() []tgs.BotCommandScope {
	return this.register
}

func (this *OPManga) Targets() *Targets {
	return this.targets
}

func (this *OPManga) AddCommands(c *tgs.MyBotCommands) {
	c.Add(BotCommandOPManga+"list", "List all available chapters", this.Register())
	c.Add(BotCommandOPManga, "Request a chapter", this.Register())
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

	this.reply <- &Reply{
		Message:  &msg,
		Timeout:  time.Minute * 5,
		Callback: this.replyCallback,
	}

	return nil
}

func (this *OPManga) arcs() ([]OPMangaArc, error) {
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

	root, err := os.ReadDir(this.path)
	if err != nil {
		return nil, err
	}

	arcs := make([]OPMangaArc, 0)

	for _, e1 := range root {
		if !e1.IsDir() {
			continue // Just ignore all non directories
		}

		sub, err := os.ReadDir(filepath.Join(this.path, e1.Name()))
		if err != nil {
			continue // Ignore for now
		}

		arc := OPMangaArc{
			Chapters: []OPMangaChapter{},
		}

		if sp := strings.SplitN(e1.Name(), " ", 2); len(sp) < 2 {
			arc.Name = e1.Name()
		} else {
			arc.Name = sp[1] // Ignore the prefixed number (ex.: "016 Thousand Sunny Arc")
		}

		for _, e2 := range sub {
			if e2.IsDir() {
				continue // Skip all directories
			}

			if filepath.Ext(e2.Name()) != ".pdf" {
				continue // Allow only pdf
			}

			chapter := OPMangaChapter{
				Path: filepath.Join(this.path, e1.Name(), e2.Name()),
			}

			// Parse file: "0441 Duell auf Banaro Island.pdf", remove ".pdf", Get prefixed
			// number and chapter name
			fileSplit := strings.SplitN(strings.TrimSuffix(e2.Name(), ".pdf"), " ", 2)

			if n, err := strconv.Atoi(fileSplit[0]); err != nil {
				return nil, err
			} else {
				chapter.Number = n
			}

			chapter.Name = fileSplit[1]

			arc.Chapters = append(arc.Chapters, chapter)
		}

		arcs = append(arcs, arc)
	}

	return arcs, nil
}

func (this *OPManga) isListCommand(c string) bool {
	return c == BotCommandOPManga[1:]+"list"
}

func (this *OPManga) handleListCommand(m *tgbotapi.Message) error {
	arcs, err := this.arcs()
	if err != nil {
		return err
	}

	content, err := GetTemplateData(&OPMangaTemplateData{
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
		"command", BotCommandOPManga,
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
			if c.Number == n {
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

	return fmt.Errorf("chapter numbrer %d not found", n)
}
