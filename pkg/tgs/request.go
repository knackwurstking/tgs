package tgs

import (
	"encoding/json"
)

type RequestGetMe struct {
	API
}

func NewRequestGetMe(api API) *RequestGetMe {
	return &RequestGetMe{
		API: api,
	}
}

func (r *RequestGetMe) Command() Command {
	return CommandGetMe
}

func (r *RequestGetMe) Send() (*ResponseGetMe, error) {
	data, err := r.API.Send(r)
	if err != nil {
		return nil, err
	}

	var response ResponseGetMe
	return &response, json.Unmarshal(data, &response)
}

type RequestGetUpdates struct {
	API

	Offset         *int     `json:"offset" yaml:"offset"`
	Limit          *int     `json:"limit" yaml:"limit"`     // Limit defaults to 100
	Timeout        *int     `json:"timeout" yaml:"timeout"` // Timeout defaults to 0 (Short Polling)
	AllowedUpdates []string `json:"allowed_updates" yaml:"allowed_updates"`
}

func NewRequestGetUpdates(api API) *RequestGetUpdates {
	return &RequestGetUpdates{
		API: api,
	}
}

func (r *RequestGetUpdates) Command() Command {
	return CommandGetUpdates
}

func (r *RequestGetUpdates) Send() (*ResponseGetUpdates, error) {
	data, err := r.API.Send(r)
	if err != nil {
		return nil, err
	}

	var response ResponseGetUpdates
	return &response, json.Unmarshal(data, &response)
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
