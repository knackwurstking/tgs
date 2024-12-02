package botcommand

import (
	"errors"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	TimeoutError = errors.New("TimeoutError")
)

type Reply struct {
	Message  *tgbotapi.Message
	Timeout  time.Duration
	Callback func(message *tgbotapi.Message) error

	done chan error
}

func (this *Reply) Run(message *tgbotapi.Message) {
	if this.Callback == nil {
		return
	}

	if err := this.Callback(message); err != nil {
		this.Done() <- err
	} else {
		this.Done() <- nil
	}
}

func (this *Reply) StartTimeout() {
	defer func() {
		if r := recover(); r != nil {
			slog.Debug("Recovered", "recover", recover())
		}
	}()

	if this.done == nil {
		this.done = make(chan error)
	}

	time.Sleep(this.Timeout)
	this.done <- TimeoutError
}

func (this *Reply) Done() chan error {
	if this.done == nil {
		this.done = make(chan error)
	}

	return this.done
}

func (this *Reply) Close() {
	defer recover()
	close(this.done)
}
