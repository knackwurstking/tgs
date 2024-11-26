package tgs

import (
	"encoding/json"
	"errors"

	"github.com/knackwurstking/tgs/pkg/data"
)

var (
	MissingAPIError = errors.New("missing API")
)

type RequestGetMe struct {
	API
}

func (*RequestGetMe) Command() Command {
	return CommandGetMe
}

func (this *RequestGetMe) Send() (*ResponseGetMe, error) {
	if this.API == nil {
		return nil, MissingAPIError
	}

	data, err := this.API.SendRequest(this)
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

func (*RequestGetUpdates) Command() Command {
	return CommandGetUpdates
}

func (this *RequestGetUpdates) Send() (*ResponseGetUpdates, error) {
	if this.API == nil {
		return nil, MissingAPIError
	}

	data, err := this.API.SendRequest(this)
	if err != nil {
		return nil, err
	}

	var response ResponseGetUpdates
	return &response, json.Unmarshal(data, &response)
}

type RequestSetMyCommands struct {
	API

	Commands     []data.BotCommand    `json:"commands"`      // A JSON-serialized list of bot commands to be set as the list of the bot's commands. At most 100 commands can be specified.
	Scope        data.BotCommandScope `json:"scope"`         // [Optional] A JSON-serialized object, describing scope of users for which the commands are relevant. Defaults to BotCommandScopeDefault.
	LanguageCode string               `json:"language_code"` // [Optional] A two-letter ISO 639-1 language code. If empty, commands will be applied to all users from the given scope, for whose language there are no dedicated commands
}

func (*RequestSetMyCommands) Command() Command {
	return CommandSetMyCommands
}

func (this *RequestSetMyCommands) Send() (*ResponseSetMyCommands, error) {
	if this.API == nil {
		return nil, MissingAPIError
	}

	data, err := this.API.SendRequest(this)
	if err != nil {
		return nil, err
	}

	var response ResponseSetMyCommands
	return &response, json.Unmarshal(data, &response)
}

type RequestSendMessage struct {
	API

	// NOTE: Valid types: int or string, only use int for now
	ChatID int `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)

	Text string `json:"text"` // Text of the message to be sent, 1-4096 characters after entities parsing

	BusinessConnectionID string                  `json:"business_connection_id"` // [Optional] Unique identifier of the business connection on behalf of which the message will be sent
	MessageThreadID      int                     `json:"message_thread_id"`      // [Optional] Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	ParseMode            string                  `json:"parse_mode"`             // [Optional] Mode for parsing entities in the message text. See formatting options for more details.
	Entities             []data.MessageEntity    `json:"entities"`               // [Optional] A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode
	LinkPreviewOptions   data.LinkPreviewOptions `json:"link_preview_options"`   // [Optional] Link preview generation options for the message
	DisableNotification  bool                    `json:"disable_notification"`   // [Optional] Sends the message silently. Users will receive a notification with no sound.
	ProtectContent       bool                    `json:"protect_content"`        // [Optional] Protects the contents of the sent message from forwarding and saving
	AllowPaidBroadcast   bool                    `json:"allow_paid_broadcast"`   // [Optional] Pass True to allow up to 1000 messages per second, ignoring broadcasting limits for a fee of 0.1 Telegram Stars per message. The relevant Stars will be withdrawn from the bot's balance
	MessageEffectID      string                  `json:"message_effect_id"`      // [Optional] Unique identifier of the message effect to be added to the message; for private chats only
	ReplyParameters      data.ReplyParameters    `json:"reply_parameters"`       // [Optional] Description of the message to reply to

	// NOTE: Ignore this for now
	//ReplyMarkup      (InlineKeyboardMarkup or ReplyKeyboardMarkup or ReplyKeyboardRemove or ForceReply)    `json:"reply_markup"`   // [Optional] Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove a reply keyboard or to force a reply from the user
}

func (*RequestSendMessage) Command() Command {
	return CommandSendMessage
}

func (this *RequestSendMessage) Send() (*ResponseSendMessage, error) {
	if this.API == nil {
		return nil, MissingAPIError
	}

	data, err := this.API.SendRequest(this)
	if err != nil {
		return nil, err
	}

	var response ResponseSendMessage
	return &response, json.Unmarshal(data, &response)
}
