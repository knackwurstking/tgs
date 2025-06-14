package tgs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReplyCallbacks map[int]func(message *tgbotapi.Message) error

func (r ReplyCallbacks) Set(id int, cb func(message *tgbotapi.Message) error) {
	r[id] = cb
}

func (r ReplyCallbacks) Get(id int) (cb func(message *tgbotapi.Message) error, ok bool) {
	cb, ok = r[id]
	return cb, ok
}

func (r ReplyCallbacks) Delete(id int) {
	delete(r, id)
}
