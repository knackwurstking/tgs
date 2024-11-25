package data

type Message struct {
	MessageID int  `json:"message_id"`
	Date      int  `json:"date"`
	Chat      Chat `json:"chat"`

	MessageThreadID    int             `json:"message_thread_id,omitempty"`    // [Optional]
	From               User            `json:"from,omitempty"`                 // [Optional]
	SenderChat         *Chat           `json:"sender_chat,omitempty"`          // [Optional]
	IsTopicMessage     *bool           `json:"is_topic_message,omitempty"`     // [Optional]
	IsAutomaticForward *bool           `json:"is_automatic_forward,omitempty"` // [Optional]
	ReplyToMessage     *Message        `json:"reply_to_message,omitempty"`     // [Optional]
	ViaBot             *User           `json:"via_bot,omitempty"`              // [Optional]
	EditDate           int             `json:"edit_date,omitempty"`            // [Optional]
	IsFromOffline      *bool           `json:"is_from_offline,omitempty"`      // [Optional]
	Text               string          `json:"text,omitempty"`                 // [Optional]
	Entities           []MessageEntity `json:"entities,omitempty"`             // [Optional]
}
