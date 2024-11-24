package data

type User struct {
	ID                      int     `json:"id" yaml:"id"`
	IsBot                   bool    `json:"is_bot" yaml:"is_bot"`
	FirstName               string  `json:"first_name" yaml:"first_name"`
	LastName                *string `json:"last_name" yaml:"last_name"`
	Username                *string `json:"username" yaml:"username"`
	LanguageCode            *string `json:"language_code" yaml:"language_code"`
	IsPremium               *bool   `json:"is_premium" yaml:"is_premium"`
	AddedToAttachmentMenu   *bool   `json:"added_to_attachment_menu" yaml:"added_to_attachment_menu"`
	CanJoinGroups           *bool   `json:"can_join_groups" yaml:"can_join_groups"`
	CanReadAllGroupMessages *bool   `json:"can_read_all_group_messages" yaml:"can_read_all_group_messages"`
	SupportsInlineQueries   *bool   `json:"supports_inline_queries" yaml:"supports_inline_queries"`
	CanConnectToBusiness    *bool   `json:"can_connect_to_business" yaml:"can_connect_to_business"`
	HasMainWebApp           *bool   `json:"has_main_web_app" yaml:"has_main_web_app"`
}

type Chat struct {
	ID        int     `json:"id" yaml:"id"`                 // Unique identifier for this chat. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	Type      string  `json:"type" yaml:"type"`             // Possible types: "private", “group”, “supergroup” or “channel”
	Title     *string `json:"title" yaml:"title"`           // Optional. Title, for supergroups, channels and group chats
	Username  *string `json:"username" yaml:"username"`     // Optional. Username, for private chats, supergroups and channels if available
	FirstName *string `json:"first_name" yaml:"first_name"` // Optional. First name of the other party in a private chat
	LastName  *string `json:"last_name" yaml:"last_name"`   // Optional. Last name of the other party in a private chat
	IsForum   *bool   `json:"is_forum" yaml:"is_forum"`     // Optional. True, if the supergroup chat is a forum (has topics enabled)
}

type MessageOrigin struct {
	// TODO: ...
}

type ExternalReplyInfo struct {
	// TODO: ...
}

type TextQuote struct {
	// TODO: ...
}

type Story struct {
	// TODO: ...
}

type MessageEntity struct {
	Type          string  `json:"type" yaml:"type"`
	Offset        int     `json:"offset" yaml:"offset"`
	Length        int     `json:"length" yaml:"length"`
	URL           *string `json:"url" yaml:"url"`
	User          *User   `json:"user" yaml:"user"`
	Language      *string `json:"language" yaml:"language"`
	CustomEmojiID *string `json:"custom_emoji_id" yaml:"custom_emoji_id"`
}

type LinkPreviewOptions struct {
	// TODO: ...
}

type Animation struct {
	// TODO: ...
}

type Audio struct {
	// TODO: ...
}

type Document struct {
	// TODO: ...
}

type PaidMediaInfo struct {
	// TODO: ...
}

type PhotoSize struct {
	// TODO: ...
}

type Sticker struct {
	// TODO: ...
}

type Video struct {
	// TODO: ...
}

type VideoNote struct {
	// TODO: ...
}

type Voice struct {
	// TODO: ...
}

type Caption struct {
	// TODO: ...
}

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

type BusinessConnection struct {
	// TODO: ...
}

type BusinessMessagesDeleted struct {
	// TODO: ...
}

type MessageReactionUpdated struct {
	// TODO: ...
}

type MessageReactionCountUpdated struct {
	// TODO: ...
}

type InlineQuery struct {
	// TODO: ...
}

type ChosenInlineResult struct {
	// TODO: ...
}

type CallbackQuery struct {
	// TODO: ...
}

type ShippingQuery struct {
	// TODO: ...
}

type PreCheckoutQuery struct {
	// TODO: ...
}

type PaidMediaPurchased struct {
	// TODO: ...
}

type Poll struct {
	// TODO: ...
}

type PollAnswer struct {
	// TODO: ...
}

type ChatMemberUpdated struct {
	// TODO: ...
}

type ChatJoinRequest struct {
	// TODO: ...
}

type ChatBoostUpdated struct {
	// TODO: ...
}

type ChatBoostRemoved struct {
	// TODO: ...
}
