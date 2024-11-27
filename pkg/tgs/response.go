package tgs

import (
	"github.com/knackwurstking/tgs/pkg/data"
)

type ResponseGetMe struct {
	OK          bool      `json:"ok" yaml:"ok"`
	Result      data.User `json:"result,omitempty" yaml:"result,omitempty"`
	ErrorCode   int       `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
}

type ResponseGetUpdates struct {
	OK          bool          `json:"ok" yaml:"ok"`
	Result      []data.Update `json:"result" yaml:"result"`
	ErrorCode   int           `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
}

type ResponseSetMyCommands struct {
	OK          bool   `json:"ok" yaml:"ok"`
	Result      bool   `json:"result" yaml:"result"`
	ErrorCode   int    `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type ResponseDeleteMyCommands struct {
	OK          bool   `json:"ok" yaml:"ok"`
	Result      bool   `json:"result" yaml:"result"`
	ErrorCode   int    `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type ResponseSendMessage struct {
	OK          bool         `json:"ok" yaml:"ok"`
	Result      data.Message `json:"result" yaml:"result"`
	ErrorCode   int          `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
}
