package tgs

import (
	"errors"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReplyCallback struct {
	ID int // ID should be the reply message id
	Fn func(message *tgbotapi.Message)
}

func NewReplyCallback(id int, fn func(message *tgbotapi.Message)) *ReplyCallback {
	return &ReplyCallback{
		ID: id,
		Fn: fn,
	}
}

type ReplyCallbacks []ReplyCallback

func (r ReplyCallbacks) Set(callbacks ...*ReplyCallback) {
	// TODO: ...
}

func (r ReplyCallbacks) Add(callbacks ...*ReplyCallback) {
	// TODO: ...
}

func (r ReplyCallbacks) Get(id int) *ReplyCallback {
	// TODO: ...

	return nil
}

func (r *ReplyCallbacks) Delete(id int) {
	// TODO: ...
}

// NOTE: Old stuff below

var ErrorReplyTimeout = errors.New("TimeoutError")

type Reply struct {
	Message  *tgbotapi.Message
	Callback func(message *tgbotapi.Message) error
	done     chan error
	Timeout  time.Duration
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
	this.done <- ErrorReplyTimeout
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
