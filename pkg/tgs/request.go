package tgs

import (
	"encoding/json"
)

type Request interface {
	Command() Command
}

type RequestGetMe struct {
	TelegramBotAPI
}

func NewRequestGetMe(api TelegramBotAPI) *RequestGetMe {
	return &RequestGetMe{
		TelegramBotAPI: api,
	}
}

func (r *RequestGetMe) Command() Command {
	return CommandGetMe
}

func (r *RequestGetMe) Send() (*ResponseGetMe, error) {
	data, err := r.TelegramBotAPI.SendRequest(r)
	if err != nil {
		return nil, err
	}

	var response ResponseGetMe
	return &response, json.Unmarshal(data, &response)
}

type RequestGetUpdates struct {
	TelegramBotAPI

	Offset         *int     `json:"offset" yaml:"offset"`
	Limit          *int     `json:"limit" yaml:"limit"`     // Limit defaults to 100
	Timeout        *int     `json:"timeout" yaml:"timeout"` // Timeout defaults to 0 (Short Polling)
	AllowedUpdates []string `json:"allowed_updates" yaml:"allowed_updates"`
}

func NewRequestGetUpdates(api TelegramBotAPI) *RequestGetUpdates {
	return &RequestGetUpdates{
		TelegramBotAPI: api,
	}
}

func (r *RequestGetUpdates) Command() Command {
	return CommandGetUpdates
}

func (r *RequestGetUpdates) Send() (*ResponseGetUpdates, error) {
	data, err := r.TelegramBotAPI.SendRequest(r)
	if err != nil {
		return nil, err
	}

	var response ResponseGetUpdates
	return &response, json.Unmarshal(data, &response)
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
