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

type Message struct {
	MessageID            int                `json:"message_id"`
	MessageThreadID      *int               `json:"message_thread_id"`
	From                 *User              `json:"from"`
	SenderChat           *Chat              `json:"sender_chat"`
	SenderBoostCount     *int               `json:"sender_boost_count"`
	SenderBusinessBot    *User              `json:"sender_business_bot"`
	Date                 int                `json:"date"`
	BusinessConnectionID *int               `json:"business_connection_id"`
	Chat                 Chat               `json:"chat"`
	ForwardOrigin        *MessageOrigin     `json:"forward_origin"`
	IsTopicMessage       *bool              `json:"is_topic_message"`
	IsAutomaticForward   *bool              `json:"is_automatic_forward"`
	ReplyToMessage       *Message           `json:"reply_to_message"`
	ExternalReply        *ExternalReplyInfo `json:"external_reply"`
	Quote                *TextQuote         `json:"quote"`
	ReplyToStory         *Story             `json:"reply_to_story"`
	// TODO: ...
}

type Update struct {
	UpdateID      int      `json:"update_id"`
	Message       *Message `json:"message"`
	EditedMessage *Message `json:"edited_message"`
	// TODO: ...
}

type ResponseGetMe User

type ResponseGetUpdates []Update

// TODO: "setMyCommands"
// TODO: "sendMessage"
