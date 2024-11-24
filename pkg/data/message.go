package data

type Message struct {
	MessageID             int                 `json:"message_id" yaml:"message_id"`
	MessageThreadID       *int                `json:"message_thread_id" yaml:"message_thread_id"`
	From                  *User               `json:"from" yaml:"from"`
	SenderChat            *Chat               `json:"sender_chat" yaml:"sender_chat"`
	SenderBoostCount      *int                `json:"sender_boost_count" yaml:"sender_boost_count"`
	SenderBusinessBot     *User               `json:"sender_business_bot" yaml:"sender_business_bot"`
	Date                  int                 `json:"date" yaml:"date"`
	BusinessConnectionID  *int                `json:"business_connection_id" yaml:"business_connection_id"`
	Chat                  Chat                `json:"chat" yaml:"chat"`
	ForwardOrigin         *MessageOrigin      `json:"forward_origin" yaml:"forward_origin"`
	IsTopicMessage        *bool               `json:"is_topic_message" yaml:"is_topic_message"`
	IsAutomaticForward    *bool               `json:"is_automatic_forward" yaml:"is_automatic_forward"`
	ReplyToMessage        *Message            `json:"reply_to_message" yaml:"reply_to_message"`
	ExternalReply         *ExternalReplyInfo  `json:"external_reply" yaml:"external_reply"`
	Quote                 *TextQuote          `json:"quote" yaml:"quote"`
	ReplyToStory          *Story              `json:"reply_to_story" yaml:"reply_to_story"`
	ViaBot                *User               `json:"via_bot" yaml:"via_bot"`
	EditDate              *int                `json:"edit_date" yaml:"edit_date"`
	HasProtectedContent   *bool               `json:"has_protected_content" yaml:"has_protected_content"`
	IsFromOffline         *bool               `json:"is_from_offline" yaml:"is_from_offline"`
	MediaGroupID          *string             `json:"media_group_id" yaml:"media_group_id"`
	AuthorSignature       *string             `json:"author_signature" yaml:"author_signature"`
	Text                  *string             `json:"text" yaml:"text"`
	Entities              []MessageEntity     `json:"entities" yaml:"entities"`
	LinkPreviewOptions    *LinkPreviewOptions `json:"link_preview_options" yaml:"link_preview_options"`
	EffectID              *string             `json:"effect_id" yaml:"effect_id"`
	Animation             *Animation          `json:"animation" yaml:"animation"`
	Audio                 *Audio              `json:"audio" yaml:"audio"`
	Document              *Document           `json:"document" yaml:"document"`
	PaidMedia             *PaidMediaInfo      `json:"paid_media" yaml:"paid_media"`
	Photo                 []PhotoSize         `json:"photo" yaml:"photo"`
	Sticker               *Sticker            `json:"sticker" yaml:"sticker"`
	Story                 *Story              `json:"story" yaml:"story"`
	Video                 *Video              `json:"video" yaml:"video"`
	VideoNote             *VideoNote          `json:"video_note" yaml:"video_note"`
	Voice                 *Voice              `json:"voice" yaml:"voice"`
	Caption               *string             `json:"caption" yaml:"caption"`
	CaptionEntities       []MessageEntity     `json:"caption_entities" yaml:"caption_entities"`
	ShowCaptionAboveMedia *bool               `json:"show_caption_above_media" yaml:"show_caption_above_media"`
	HasMediaSpoiler       *bool               `json:"has_media_spoiler" yaml:"has_media_spoiler"`
	// TODO: ...
}
