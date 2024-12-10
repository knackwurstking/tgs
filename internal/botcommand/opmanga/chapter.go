package opmanga

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/knackwurstking/tgs/internal/botcommand"
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

func (chapter *Chapter) PDF() (pdf botcommand.File, err error) {
	slog.Debug("Read pdf", "path", chapter.Path)
	if data, err := os.ReadFile(chapter.Path); err != nil {
		return nil, err
	} else {
		return NewPDF(filepath.Base(chapter.Path), data), nil
	}
}
