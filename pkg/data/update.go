package data

type Update struct {
	UpdateID                int                          `json:"update_id" yaml:"update_id"`
	Message                 *Message                     `json:"message" yaml:"message"`
	EditedMessage           *Message                     `json:"edited_message" yaml:"edited_message"`
	ChannelPost             *Message                     `json:"channel_post" yaml:"channel_post"`
	EditedChannelPost       *Message                     `json:"edited_channel_post" yaml:"edited_channel_post"`
	BusinessConnection      *BusinessConnection          `json:"business_connection" yaml:"business_connection"`
	BusinessMessage         *Message                     `json:"business_message" yaml:"business_message"`
	EditedBusinessMessage   *Message                     `json:"edited_business_message" yaml:"edited_business_message"`
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages" yaml:"deleted_business_messages"`
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction" yaml:"message_reaction"`
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count" yaml:"message_reaction_count"`
	InlineQuery             *InlineQuery                 `json:"inline_query" yaml:"inline_query"`
	ChosenInlineResult      *ChosenInlineResult          `json:"chosen_inline_result" yaml:"chosen_inline_result"`
	CallbackQuery           *CallbackQuery               `json:"callback_query" yaml:"callback_query"`
	ShippingQuery           *ShippingQuery               `json:"shipping_query" yaml:"shipping_query"`
	PreCheckoutQuery        *PreCheckoutQuery            `json:"pre_checkout_query" yaml:"pre_checkout_query"`
	PurchasedPaidMedia      *PaidMediaPurchased          `json:"purchased_paid_media" yaml:"purchased_paid_media"`
	Poll                    *Poll                        `json:"poll" yaml:"poll"`
	PollAnswer              *PollAnswer                  `json:"poll_answer" yaml:"poll_answer"`
	MyChatMember            *ChatMemberUpdated           `json:"my_chat_member" yaml:"my_chat_member"`
	ChatMember              *ChatMemberUpdated           `json:"chat_member" yaml:"chat_member"`
	ChatJoinRequest         *ChatJoinRequest             `json:"chat_join_request" yaml:"chat_join_request"`
	ChatBoost               *ChatBoostUpdated            `json:"chat_boost" yaml:"chat_boost"`
	RemovedChatBoost        *ChatBoostRemoved            `json:"removed_chat_boost" yaml:"removed_chat_boost"`
}
