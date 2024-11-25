package data

type Message struct {
	MessageID          int             `json:"message_id"`
	MessageThreadID    int             `json:"message_thread_id"` // [Optional]
	From               User            `json:"from"`              // [Optional]
	SenderChat         *Chat           `json:"sender_chat"`       // [Optional]
	Date               int             `json:"date"`
	Chat               Chat            `json:"chat"`
	IsTopicMessage     bool            `json:"is_topic_message"`     // [Optional]
	IsAutomaticForward bool            `json:"is_automatic_forward"` // [Optional]
	ReplyToMessage     *Message        `json:"reply_to_message"`     // [Optional]
	ViaBot             *User           `json:"via_bot"`              // [Optional]
	EditDate           int             `json:"edit_date"`            // [Optional]
	IsFromOffline      bool            `json:"is_from_offline"`      // [Optional]
	Text               string          `json:"text"`                 // [Optional]
	Entities           []MessageEntity `json:"entities"`             // [Optional]
}
