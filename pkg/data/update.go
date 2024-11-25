package data

type Update struct {
	UpdateID                int                          `json:"update_id"`
	Message                 *Message                     `json:"message"`                   // [Optional]
	EditedMessage           *Message                     `json:"edited_message"`            // [Optional]
	ChannelPost             *Message                     `json:"channel_post"`              // [Optional]
	EditedChannelPost       *Message                     `json:"edited_channel_post"`       // [Optional]
	BusinessConnection      *BusinessConnection          `json:"business_connection"`       // [Optional]
	BusinessMessage         *Message                     `json:"business_message"`          // [Optional]
	EditedBusinessMessage   *Message                     `json:"edited_business_message"`   // [Optional]
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages"` // [Optional]
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction"`          // [Optional]
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count"`    // [Optional]
	InlineQuery             *InlineQuery                 `json:"inline_query"`              // [Optional]
	ChosenInlineResult      *ChosenInlineResult          `json:"chosen_inline_result"`      // [Optional]
	CallbackQuery           *CallbackQuery               `json:"callback_query"`            // [Optional]
	ShippingQuery           *ShippingQuery               `json:"shipping_query"`            // [Optional]
	PreCheckoutQuery        *PreCheckoutQuery            `json:"pre_checkout_query"`        // [Optional]
	PurchasedPaidMedia      *PaidMediaPurchased          `json:"purchased_paid_media"`      // [Optional]
	Poll                    *Poll                        `json:"poll"`                      // [Optional]
	PollAnswer              *PollAnswer                  `json:"poll_answer"`               // [Optional]
	MyChatMember            *ChatMemberUpdated           `json:"my_chat_member"`            // [Optional]
	ChatMember              *ChatMemberUpdated           `json:"chat_member"`               // [Optional]
	ChatJoinRequest         *ChatJoinRequest             `json:"chat_join_request"`         // [Optional]
	ChatBoost               *ChatBoostUpdated            `json:"chat_boost"`                // [Optional]
	RemovedChatBoost        *ChatBoostRemoved            `json:"removed_chat_boost"`        // [Optional]
}
