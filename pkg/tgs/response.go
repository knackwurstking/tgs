package tgs

import "github.com/knackwurstking/tgs/pkg/data"

type ResponseGetMe struct {
	OK     bool      `json:"ok" yaml:"ok"`
	Result data.User `json:"result" yaml:"result"`
}

type ResponseGetUpdates struct {
	OK     bool          `json:"ok" yaml:"ok"`
	Result []data.Update `json:"result" yaml:"result"`
}

type ResponseSetMyCommands struct {
	OK     bool `json:"ok" yaml:"ok"`
	Result bool `json:"result" yaml:"result"`
}

type ResponseDeleteMyCommands struct {
	OK     bool `json:"ok" yaml:"ok"`
	Result bool `json:"result" yaml:"result"`
}

type ResponseSendMessage struct {
	OK     bool         `json:"ok" yaml:"ok"`
	Result data.Message `json:"result" yaml:"result"`
}
