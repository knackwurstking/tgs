package tgs

type User struct {
	ID                      int     `json:"id"`
	IsBot                   bool    `json:"is_bot"`
	FirstName               string  `json:"first_name"`
	LastName                *string `json:"last_name"`
	Username                *string `json:"username"`
	LanguageCode            *string `json:"language_code"`
	IsPremium               *bool   `json:"is_premium"`
	AddedToAttachmentMenu   *bool   `json:"added_to_attachment_menu"`
	CanJoinGroups           *bool   `json:"can_join_groups"`
	CanReadAllGroupMessages *bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   *bool   `json:"supports_inline_queries"`
	CanConnectToBusiness    *bool   `json:"can_connect_to_business"`
	HasMainWebApp           *bool   `json:"has_main_web_app"`
}

type Chat struct {
	// TODO: ...
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
	Type          string  `json:"type"`
	Offset        int     `json:"offset"`
	Length        int     `json:"length"`
	URL           *string `json:"url"`
	User          *User   `json:"user"`
	Language      *string `json:"language"`
	CustomEmojiID *string `json:"custom_emoji_id"`
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
	MessageID             int                 `json:"message_id"`
	MessageThreadID       *int                `json:"message_thread_id"`
	From                  *User               `json:"from"`
	SenderChat            *Chat               `json:"sender_chat"`
	SenderBoostCount      *int                `json:"sender_boost_count"`
	SenderBusinessBot     *User               `json:"sender_business_bot"`
	Date                  int                 `json:"date"`
	BusinessConnectionID  *int                `json:"business_connection_id"`
	Chat                  Chat                `json:"chat"`
	ForwardOrigin         *MessageOrigin      `json:"forward_origin"`
	IsTopicMessage        *bool               `json:"is_topic_message"`
	IsAutomaticForward    *bool               `json:"is_automatic_forward"`
	ReplyToMessage        *Message            `json:"reply_to_message"`
	ExternalReply         *ExternalReplyInfo  `json:"external_reply"`
	Quote                 *TextQuote          `json:"quote"`
	ReplyToStory          *Story              `json:"reply_to_story"`
	ViaBot                *User               `json:"via_bot"`
	EditDate              *int                `json:"edit_date"`
	HasProtectedContent   *bool               `json:"has_protected_content"`
	IsFromOffline         *bool               `json:"is_from_offline"`
	MediaGroupID          *string             `json:"media_group_id"`
	AuthorSignature       *string             `json:"author_signature"`
	Text                  *string             `json:"text"`
	Entities              []MessageEntity     `json:"entities"`
	LinkPreviewOptions    *LinkPreviewOptions `json:"link_preview_options"`
	EffectID              *string             `json:"effect_id"`
	Animation             *Animation          `json:"animation"`
	Audio                 *Audio              `json:"audio"`
	Document              *Document           `json:"document"`
	PaidMedia             *PaidMediaInfo      `json:"paid_media"`
	Photo                 []PhotoSize         `json:"photo"`
	Sticker               *Sticker            `json:"sticker"`
	Story                 *Story              `json:"story"`
	Video                 *Video              `json:"video"`
	VideoNote             *VideoNote          `json:"video_note"`
	Voice                 *Voice              `json:"voice"`
	Caption               *string             `json:"caption"`
	CaptionEntities       []MessageEntity     `json:"caption_entities"`
	ShowCaptionAboveMedia *bool               `json:"show_caption_above_media"`
	HasMediaSpoiler       *bool               `json:"has_media_spoiler"`
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

type Update struct {
	UpdateID                int                          `json:"update_id"`
	Message                 *Message                     `json:"message"`
	EditedMessage           *Message                     `json:"edited_message"`
	ChannelPost             *Message                     `json:"channel_post"`
	EditedChannelPost       *Message                     `json:"edited_channel_post"`
	BusinessConnection      *BusinessConnection          `json:"business_connection"`
	BusinessMessage         *Message                     `json:"business_message"`
	EditedBusinessMessage   *Message                     `json:"edited_business_message"`
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages"`
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction"`
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count"`
	InlineQuery             *InlineQuery                 `json:"inline_query"`
	ChosenInlineResult      *ChosenInlineResult          `json:"chosen_inline_result"`
	CallbackQuery           *CallbackQuery               `json:"callback_query"`
	ShippingQuery           *ShippingQuery               `json:"shipping_query"`
	PreCheckoutQuery        *PreCheckoutQuery            `json:"pre_checkout_query"`
	PurchasedPaidMedia      *PaidMediaPurchased          `json:"purchased_paid_media"`
	Poll                    *Poll                        `json:"poll"`
	PollAnswer              *PollAnswer                  `json:"poll_answer"`
	MyChatMember            *ChatMemberUpdated           `json:"my_chat_member"`
	ChatMember              *ChatMemberUpdated           `json:"chat_member"`
	ChatJoinRequest         *ChatJoinRequest             `json:"chat_join_request"`
	ChatBoost               *ChatBoostUpdated            `json:"chat_boost"`
	RemovedChatBoost        *ChatBoostRemoved            `json:"removed_chat_boost"`
}

type ResponseGetMe struct {
	OK     bool `json:"ok"`
	Result User `json:"result"`
}

type ResponseGetUpdates struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// TODO: "setMyCommands"
// TODO: "sendMessage"
