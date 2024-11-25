package data

type ReplyParameters struct {
	MessageID int `json:"message_id"` // Identifier of the message that will be replied to in the current chat, or in the chat chat_id if it is specified

	// NOTE: Valid types: int or string, only use int for now
	ChatID int `json:"chat_id,omitempty"` // [Optional] If the message to be replied to is from a different chat, unique identifier for the chat or username of the channel (in the format @channelusername). Not supported for messages sent on behalf of a business account.

	AllowSendingWithoutReply *bool           `json:"allow_sending_without_reply,omitempty"` // [Optional] Pass True if the message should be sent even if the specified message to be replied to is not found. Always False for replies in another chat or forum topic. Always True for messages sent on behalf of a business account.
	Quote                    string          `json:"quote,omitempty"`                       // [Optional] Quoted part of the message to be replied to; 0-1024 characters after entities parsing. The quote must be an exact substring of the message to be replied to, including bold, italic, underline, strikethrough, spoiler, and custom_emoji entities. The message will fail to send if the quote isn't found in the original message.
	QuoteParseMode           string          `json:"quote_parse_mode,omitempty"`            // [Optional] Mode for parsing entities in the quote. See formatting options for more details.
	QuoteEntities            []MessageEntity `json:"quote_entities,omitempty"`              // [Optional] A JSON-serialized list of special entities that appear in the quote. It can be specified instead of quote_parse_mode.
	QuotePosition            int             `json:"quote_position,omitempty"`              // [Optional] Position of the quote in the original message in UTF-16 code units
}
