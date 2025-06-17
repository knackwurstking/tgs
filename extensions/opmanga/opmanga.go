package opmanga

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"

	"github.com/knackwurstking/tgs/internal/templates"
	"github.com/knackwurstking/tgs/pkg/tgs"
)

type Chapter struct {
	Path string

	name   string
	number int
}

func NewChapter(path string) (*Chapter, error) {
	c := Chapter{Path: path}

	// Parse file: "0441 Duell auf Banaro Island.pdf", remove ".pdf", Get prefixed
	// number and chapter name
	fS := strings.SplitN(strings.TrimSuffix(filepath.Base(path), ".pdf"), " ", 2)

	if n, err := strconv.Atoi(fS[0]); err != nil {
		return nil, err
	} else {
		c.number = n
	}

	c.name = fS[1]

	return &c, nil
}

func (chapter *Chapter) Name() string {
	return chapter.name
}

func (chapter *Chapter) Number() int {
	return chapter.number
}

type Arc struct {
	Name     string
	Chapters []*Chapter
}

type TemplateData struct {
	PageTitle string
	Arcs      []Arc
}

func (td *TemplateData) Patterns() []string {
	return []string{
		"data/index.go.html",
		"data/opmangalist.go.html", // block: content
		"data/ui-v2.0.0.css",       // block: ui
		"data/styles.css",          // block: styles
	}
}

type Data struct {
	Targets  *tgs.Targets          `yaml:"targets,omitempty"`
	Register []tgs.BotCommandScope `yaml:"register,omitempty"`
	Path     string                `yaml:"path"`
}

type OPManga struct {
	*tgbotapi.BotAPI

	data      *Data
	callbacks tgs.ReplyCallbacks
}

func New(api *tgbotapi.BotAPI) *OPManga {
	return &OPManga{
		BotAPI: api,
		data: &Data{
			Targets:  tgs.NewTargets(),
			Register: make([]tgs.BotCommandScope, 0),
			// Reply: ,
		},
		callbacks: tgs.ReplyCallbacks{},
	}
}

func NewExtension(api *tgbotapi.BotAPI) tgs.Extension {
	return New(api)
}

func (o *OPManga) Name() string {
	return "OPManga"
}

func (o *OPManga) SetBot(api *tgbotapi.BotAPI) {
	o.BotAPI = api
}

func (o *OPManga) ConfigPath() string {
	return "opmanga.yaml"
}

func (o *OPManga) MarshalYAML() (any, error) {
	return o.data, nil
}

func (o *OPManga) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(o.data)
}

func (o *OPManga) Commands(mbc *tgs.MyBotCommands) {
	mbc.Add("/opmanga", "Request a chapter", o.data.Register)
	mbc.Add("/opmangalist", "List all available chapters", o.data.Register)
}

func (o *OPManga) Is(message *tgbotapi.Message) bool {
	replyMessageID := message.ReplyToMessage.MessageID
	if replyMessageID != 0 {
		if _, ok := o.callbacks.Get(replyMessageID); ok {
			return true
		}
	}

	return strings.HasPrefix(message.Command(), "opmanga")
}

func (o *OPManga) Handle(message *tgbotapi.Message) error {
	if o.BotAPI == nil {
		panic("BotAPI is nil!")
	}

	if ok := tgs.CheckTargets(message, o.data.Targets); !ok {
		return errors.New("invalid target")
	}

	replyMessageID := message.ReplyToMessage.MessageID
	if replyMessageID != 0 {
		if cb, ok := o.callbacks.Get(replyMessageID); ok {
			err := cb(message)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("`error: %s`", err))
				msg.ParseMode = "MarkdownV2"
				msg.ReplyToMessageID = replyMessageID
				_, err = o.Send(msg)
			}

			return err
		}
	}

	switch command := message.Command(); command {
	case "opmanga":
		msgConfig := tgbotapi.NewMessage(
			message.Chat.ID,
			"Hey there\\! Can you send me the chapter number?\n\n"+
				"You’ll have about 5 minutes to respond to this message\\.\n\n"+
				">You need to reply to this message for this to work\\.",
		)
		msgConfig.ReplyToMessageID = message.MessageID
		msgConfig.ParseMode = "MarkdownV2"

		msg, err := o.Send(msgConfig)
		if err != nil {
			return err
		}

		o.callbacks.Set(msg.MessageID, o.replyCallbackOPMangaCommand)

		go func() { // Auto Delete Function
			time.Sleep(time.Minute * 5)
			o.callbacks.Delete(msg.MessageID)
		}()

		return nil
	case "opmangalist":
		arcs, err := o.arcs()
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

		documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
			Name:  "opmanga-chapters.html",
			Bytes: content,
		})
		documentConfig.ReplyToMessageID = message.MessageID

		_, err = o.Send(documentConfig)
		return err
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (o *OPManga) replyCallbackOPMangaCommand(message *tgbotapi.Message) error {
	slog.Debug("Handle reply callback",
		"command", message.Command(),
		"message.MessageID", message.MessageID,
		"message.Text", message.Text,
	)

	arcs, err := o.arcs()
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
				chatID := message.Chat.ID

				documentConfig := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(c.Path))
				documentConfig.ReplyToMessageID = message.ReplyToMessage.MessageID

				_, err = o.Send(documentConfig)
				if err != nil {
					slog.Error("Send chapter", "error", err)
					return err
				}

				return nil
			}
		}
	}

	return fmt.Errorf("chapter number %d not found", chapterNumber)
}

func (o *OPManga) arcs() ([]Arc, error) {
	if o.data.Path == "" {
		return nil, fmt.Errorf("missing path")
	}

	info, err := os.Stat(o.data.Path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("nope, need an directory here")
	}

	dirEntries, err := os.ReadDir(o.data.Path)
	if err != nil {
		return nil, err
	}

	arcs := make([]Arc, 0)

	for _, e := range dirEntries {
		if !e.IsDir() {
			continue // Just ignore all non directories
		}

		subDirEntries, err := os.ReadDir(filepath.Join(o.data.Path, e.Name()))
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
				filepath.Join(o.data.Path, e.Name(), e2.Name()),
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
