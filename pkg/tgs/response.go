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

// TODO: "setMyCommands"

type ResponseSendMessage struct {
	OK     bool         `json:"ok" yaml:"ok"`
	Result data.Message `json:"result" yaml:"result"`
}
