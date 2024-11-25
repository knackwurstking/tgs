package data

type Message struct {
	MessageID             int                 `json:"message_id"`
	MessageThreadID       int                 `json:"message_thread_id"`   // [Optional]
	From                  User                `json:"from"`                // [Optional]
	SenderChat            *Chat               `json:"sender_chat"`         // [Optional]
	SenderBoostCount      int                 `json:"sender_boost_count"`  // [Optional]
	SenderBusinessBot     *User               `json:"sender_business_bot"` // [Optional]
	Date                  int                 `json:"date"`
	BusinessConnectionID  int                 `json:"business_connection_id"` // [Optional]
	Chat                  Chat                `json:"chat"`
	ForwardOrigin         *MessageOrigin      `json:"forward_origin"`           // [Optional]
	IsTopicMessage        bool                `json:"is_topic_message"`         // [Optional]
	IsAutomaticForward    bool                `json:"is_automatic_forward"`     // [Optional]
	ReplyToMessage        *Message            `json:"reply_to_message"`         // [Optional]
	ExternalReply         *ExternalReplyInfo  `json:"external_reply"`           // [Optional]
	Quote                 *TextQuote          `json:"quote"`                    // [Optional]
	ReplyToStory          *Story              `json:"reply_to_story"`           // [Optional]
	ViaBot                *User               `json:"via_bot"`                  // [Optional]
	EditDate              int                 `json:"edit_date"`                // [Optional]
	HasProtectedContent   bool                `json:"has_protected_content"`    // [Optional]
	IsFromOffline         bool                `json:"is_from_offline"`          // [Optional]
	MediaGroupID          string              `json:"media_group_id"`           // [Optional]
	AuthorSignature       string              `json:"author_signature"`         // [Optional]
	Text                  string              `json:"text"`                     // [Optional]
	Entities              []MessageEntity     `json:"entities"`                 // [Optional]
	LinkPreviewOptions    *LinkPreviewOptions `json:"link_preview_options"`     // [Optional]
	EffectID              string              `json:"effect_id"`                // [Optional]
	Animation             *Animation          `json:"animation"`                // [Optional]
	Audio                 *Audio              `json:"audio"`                    // [Optional]
	Document              *Document           `json:"document"`                 // [Optional]
	PaidMedia             *PaidMediaInfo      `json:"paid_media"`               // [Optional]
	Photo                 []PhotoSize         `json:"photo"`                    // [Optional]
	Sticker               *Sticker            `json:"sticker"`                  // [Optional]
	Story                 *Story              `json:"story"`                    // [Optional]
	Video                 *Video              `json:"video"`                    // [Optional]
	VideoNote             *VideoNote          `json:"video_note"`               // [Optional]
	Voice                 *Voice              `json:"voice"`                    // [Optional]
	Caption               string              `json:"caption"`                  // [Optional]
	CaptionEntities       []MessageEntity     `json:"caption_entities"`         // [Optional]
	ShowCaptionAboveMedia bool                `json:"show_caption_above_media"` // [Optional]
	HasMediaSpoiler       bool                `json:"has_media_spoiler"`        // [Optional]

	// TODO: ...
}
