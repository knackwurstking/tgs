package data

// TODO: Keep only needed fields
type Message struct {
	MessageID             int             `json:"message_id"`
	MessageThreadID       int             `json:"message_thread_id"` // [Optional]
	From                  User            `json:"from"`              // [Optional]
	SenderChat            *Chat           `json:"sender_chat"`       // [Optional]
	Date                  int             `json:"date"`
	Chat                  Chat            `json:"chat"`
	IsTopicMessage        bool            `json:"is_topic_message"`         // [Optional]
	IsAutomaticForward    bool            `json:"is_automatic_forward"`     // [Optional]
	ReplyToMessage        *Message        `json:"reply_to_message"`         // [Optional]
	ViaBot                *User           `json:"via_bot"`                  // [Optional]
	EditDate              int             `json:"edit_date"`                // [Optional]
	HasProtectedContent   bool            `json:"has_protected_content"`    // [Optional]
	IsFromOffline         bool            `json:"is_from_offline"`          // [Optional]
	MediaGroupID          string          `json:"media_group_id"`           // [Optional]
	AuthorSignature       string          `json:"author_signature"`         // [Optional]
	Text                  string          `json:"text"`                     // [Optional]
	Entities              []MessageEntity `json:"entities"`                 // [Optional]
	EffectID              string          `json:"effect_id"`                // [Optional]
	Caption               string          `json:"caption"`                  // [Optional]
	CaptionEntities       []MessageEntity `json:"caption_entities"`         // [Optional]
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media"` // [Optional]
	HasMediaSpoiler       bool            `json:"has_media_spoiler"`        // [Optional]
}
