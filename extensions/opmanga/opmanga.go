package opmanga

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
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
		"data/ui-v4.3.0.css",       // block: ui
		"data/styles.css",          // block: styles
	}
}

type Data struct {
	Targets *tgs.Targets `yaml:"targets,omitempty"`
	Scopes  []tgs.Scope  `yaml:"scopes,omitempty"`
	Path    string       `yaml:"path"`
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
			Targets: tgs.NewTargets(),
			Scopes:  make([]tgs.Scope, 0),
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

func (o *OPManga) AddBotCommands(mbc *tgs.MyBotCommands) {
	mbc.Add("/opmanga", "Request a chapter", o.data.Scopes)
	mbc.Add("/opmangalist", "List all available chapters", o.data.Scopes)
}

func (o *OPManga) Is(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}

	if update.Message.ReplyToMessage != nil {
		if _, ok := o.callbacks.Get(update.Message.ReplyToMessage.MessageID); ok {
			return true
		}
	}

	return strings.HasPrefix(update.Message.Command(), "opmanga")
}

func (o *OPManga) Handle(update tgbotapi.Update) error {
	if o.BotAPI == nil {
		log.Error("OPManga extension BotAPI is nil")
		panic("BotAPI is nil!")
	}

	message := update.Message
	log.Debug("OPManga extension handling update",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"chat_id", message.Chat.ID,
		"command", message.Command(),
	)

	if ok := tgs.CheckTargets(message, o.data.Targets); !ok {
		log.Debug("OPManga request from unauthorized target",
			"user_id", message.From.ID,
			"chat_id", message.Chat.ID,
			"command", message.Command(),
		)
		return errors.New("invalid target")
	}

	if message.ReplyToMessage != nil {
		replyMessageID := message.ReplyToMessage.MessageID
		if cb, ok := o.callbacks.Get(replyMessageID); ok {
			log.Debug("Processing OPManga reply callback",
				"reply_to_message_id", replyMessageID,
				"user_text", message.Text,
			)

			err := cb(message)
			if err != nil {
				log.Error("OPManga reply callback failed",
					"reply_to_message_id", replyMessageID,
					"error", err,
				)

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
		log.Info("Initiating OPManga chapter request",
			"user_id", message.From.ID,
			"username", message.From.UserName,
		)

		msgConfig := tgbotapi.NewMessage(
			message.Chat.ID,
			"Hey there\\! Can you send me the chapter number?\n\n"+
				"You'll have about 5 minutes to respond to this message\\.\n\n"+
				">You need to reply to this message for this to work\\.",
		)
		msgConfig.ReplyToMessageID = message.MessageID
		msgConfig.ParseMode = "MarkdownV2"

		msg, err := o.Send(msgConfig)
		if err != nil {
			log.Error("Failed to send OPManga chapter request prompt",
				"error", err,
			)
			return err
		}

		o.callbacks.Set(msg.MessageID, o.replyCallbackOPMangaCommand)

		log.Debug("OPManga callback registered",
			"callback_message_id", msg.MessageID,
			"timeout_minutes", 5,
		)

		go func() { // Auto Delete Function
			time.Sleep(time.Minute * 5)
			o.callbacks.Delete(msg.MessageID)
			log.Debug("OPManga callback expired and removed",
				"callback_message_id", msg.MessageID,
			)
		}()

		return nil
	case "opmangalist":
		log.Info("Generating OPManga chapter list",
			"user_id", message.From.ID,
			"manga_path", o.data.Path,
		)

		arcs, err := o.arcs()
		if err != nil {
			log.Error("Failed to scan manga directory",
				"path", o.data.Path,
				"error", err,
			)
			return err
		}

		totalChapters := 0
		for _, arc := range arcs {
			totalChapters += len(arc.Chapters)
		}

		log.Debug("Manga directory scanned successfully",
			"arcs_found", len(arcs),
			"total_chapters", totalChapters,
		)

		content, err := templates.GetTemplateData(&TemplateData{
			PageTitle: "One Piece Manga",
			Arcs:      arcs,
		})
		if err != nil {
			log.Error("Failed to generate OPManga template",
				"error", err,
			)
			return err
		}

		documentConfig := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
			Name:  "opmanga-chapters.html",
			Bytes: content,
		})
		documentConfig.ReplyToMessageID = message.MessageID

		log.Debug("Sending OPManga chapter list document",
			"document_size_bytes", len(content),
			"arcs_count", len(arcs),
			"total_chapters", totalChapters,
		)

		_, err = o.Send(documentConfig)
		if err != nil {
			log.Error("Failed to send OPManga chapter list",
				"error", err,
			)
		}
		return err
	default:
		log.Warn("Unknown command in OPManga extension",
			"command", command,
			"user_id", message.From.ID,
		)
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (o *OPManga) replyCallbackOPMangaCommand(message *tgbotapi.Message) error {
	log.Info("Processing OPManga chapter request",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"requested_text", message.Text,
		"reply_to_message_id", message.ReplyToMessage.MessageID,
		"message_id", message.MessageID,
	)

	arcs, err := o.arcs()
	if err != nil {
		log.Error("Failed to scan manga directory for chapter request",
			"path", o.data.Path,
			"error", err,
		)
		return err
	}

	// Parse message and get episode string
	match := regexp.MustCompile(`[0-9]+`).FindString(message.Text)
	if match == "" {
		log.Warn("Invalid chapter request format",
			"user_text", message.Text,
			"user_id", message.From.ID,
		)
		return fmt.Errorf("nothing found, need a number here: %s", message.Text)
	}

	chapterNumber, err := strconv.Atoi(match)
	if err != nil {
		log.Error("Failed to parse chapter number",
			"match", match,
			"user_text", message.Text,
			"error", err,
		)
		return err
	}

	log.Debug("Searching for chapter",
		"chapter_number", chapterNumber,
		"available_arcs", len(arcs),
	)

	// Search arcs data for chapter
	for _, a := range arcs {
		for _, c := range a.Chapters {
			if c.Number() == chapterNumber {
				log.Info("Found requested chapter",
					"chapter_number", chapterNumber,
					"chapter_name", c.Name(),
					"arc_name", a.Name,
					"file_path", c.Path,
				)

				chatID := message.Chat.ID

				// Get file info for logging
				if fileInfo, err := os.Stat(c.Path); err == nil {
					log.Debug("Sending chapter file",
						"file_size_bytes", fileInfo.Size(),
						"file_name", fileInfo.Name(),
					)
				}

				documentConfig := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(c.Path))
				documentConfig.ReplyToMessageID = message.ReplyToMessage.MessageID

				_, err = o.Send(documentConfig)
				if err != nil {
					log.Error("Failed to send chapter file",
						"chapter_number", chapterNumber,
						"file_path", c.Path,
						"error", err,
					)
					return err
				}

				log.Info("Chapter sent successfully",
					"chapter_number", chapterNumber,
					"user_id", message.From.ID,
					"file_path", c.Path,
				)
				return nil
			}
		}
	}

	log.Warn("Chapter not found",
		"requested_chapter", chapterNumber,
		"user_id", message.From.ID,
		"available_arcs", len(arcs),
	)
	return fmt.Errorf("chapter number %d not found", chapterNumber)
}

func (o *OPManga) arcs() ([]Arc, error) {
	if o.data.Path == "" {
		log.Error("OPManga path not configured")
		return nil, fmt.Errorf("missing path")
	}

	log.Debug("Scanning manga directory", "path", o.data.Path)

	info, err := os.Stat(o.data.Path)
	if err != nil {
		log.Error("Failed to access manga directory",
			"path", o.data.Path,
			"error", err,
		)
		return nil, err
	}
	if !info.IsDir() {
		log.Error("Manga path is not a directory",
			"path", o.data.Path,
		)
		return nil, fmt.Errorf("nope, need an directory here")
	}

	dirEntries, err := os.ReadDir(o.data.Path)
	if err != nil {
		log.Error("Failed to read manga directory",
			"path", o.data.Path,
			"error", err,
		)
		return nil, err
	}

	log.Debug("Found directory entries",
		"path", o.data.Path,
		"entry_count", len(dirEntries),
	)

	arcs := make([]Arc, 0)

	for _, e := range dirEntries {
		if !e.IsDir() {
			log.Debug("Skipping non-directory entry", "name", e.Name())
			continue // Just ignore all non directories
		}

		log.Debug("Processing arc directory", "name", e.Name())

		subDirPath := filepath.Join(o.data.Path, e.Name())
		subDirEntries, err := os.ReadDir(subDirPath)
		if err != nil {
			log.Warn("Failed to read arc directory, skipping",
				"arc_dir", e.Name(),
				"path", subDirPath,
				"error", err,
			)
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

		chapterCount := 0
		for _, e2 := range subDirEntries {
			if e2.IsDir() {
				log.Debug("Skipping subdirectory in arc", "name", e2.Name(), "arc", arc.Name)
				continue // Skip all directories
			}

			if filepath.Ext(e2.Name()) != ".pdf" {
				log.Debug("Skipping non-PDF file", "name", e2.Name(), "arc", arc.Name)
				continue // Allow only pdf
			}

			chapterPath := filepath.Join(o.data.Path, e.Name(), e2.Name())
			chapter, err := NewChapter(chapterPath)
			if err != nil {
				log.Error("Failed to parse chapter file",
					"file_path", chapterPath,
					"arc", arc.Name,
					"error", err,
				)
				return nil, err
			}

			arc.Chapters = append(arc.Chapters, chapter)
			chapterCount++
		}

		log.Debug("Arc processed",
			"arc_name", arc.Name,
			"chapters_found", chapterCount,
		)

		arcs = append(arcs, arc)
	}

	totalChapters := 0
	for _, arc := range arcs {
		totalChapters += len(arc.Chapters)
	}

	log.Info("Manga directory scan completed",
		"path", o.data.Path,
		"arcs_found", len(arcs),
		"total_chapters", totalChapters,
	)

	return arcs, nil
}
