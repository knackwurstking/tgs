package tgs

import (
	"encoding/json"
	"fmt"
)

type RequestGetMe struct {
	API
}

func (r *RequestGetMe) Command() Command {
	return CommandGetMe
}

func (r *RequestGetMe) Send() (*ResponseGetMe, error) {
	if r.API == nil {
		return nil, fmt.Errorf("missing API")
	}

	data, err := r.API.SendRequest(r)
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

func (r *RequestGetUpdates) Command() Command {
	return CommandGetUpdates
}

func (r *RequestGetUpdates) Send() (*ResponseGetUpdates, error) {
	if r.API == nil {
		return nil, fmt.Errorf("missing API")
	}

	data, err := r.API.SendRequest(r)
	if err != nil {
		return nil, err
	}

	var response ResponseGetUpdates
	return &response, json.Unmarshal(data, &response)
}

// TODO: "setMyCommands"

type RequestSendMessage struct {
	API

	// TODO: Request data missing here for "sendMessage"
}

func (r *RequestSendMessage) Command() Command {
	return CommandSendMessage
}

// TODO: Message type missing here
func (r *RequestSendMessage) Send(message any) (*ResponseSendMessage, error) {
	if r.API == nil {
		return nil, fmt.Errorf("missing API")
	}

	data, err := r.API.SendRequest(r)
	if err != nil {
		return nil, err
	}

	var response ResponseSendMessage
	return &response, json.Unmarshal(data, &response)
}
