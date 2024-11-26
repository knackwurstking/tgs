package data

import "fmt"

const (
	BotCommandScopeTypeDefault               = BotCommandScopeType("default")
	BotCommandScopeTypeAllPrivateChats       = BotCommandScopeType("all_private_chats")
	BotCommandScopeTypeAllGroupChats         = BotCommandScopeType("all_group_chats")
	BotCommandScopeTypeAllChatAdministrators = BotCommandScopeType("all_chat_administrators")
	BotCommandScopeTypeChat                  = BotCommandScopeType("chat")
	BotCommandScopeTypeChatAdministrators    = BotCommandScopeType("chat_administrators")
	BotCommandScopeTypeChatMember            = BotCommandScopeType("chat_member")
)

type BotCommandScopeType string

type BotCommandScope struct {
	Type string `json:"type"`

	ChatID int `json:"chat_id,omitempty"` // [Optional] Types: "chat", "chat_administrators", "chat_member"
	UserID int `json:"user_id,omitempty"` // [Optional] Types: "chat_member"
}

func (this *BotCommandScope) Scope() BotCommandScopeType {
	return BotCommandScopeType(this.Type)
}

func (this *BotCommandScope) SetScope(scope BotCommandScopeType, chatID int, userID int) error {
	switch scope {
	case
		BotCommandScopeTypeDefault,
		BotCommandScopeTypeAllPrivateChats,
		BotCommandScopeTypeAllGroupChats,
		BotCommandScopeTypeAllChatAdministrators:

		this.Type = string(scope)
		break

	case
		BotCommandScopeTypeChat,
		BotCommandScopeTypeChatAdministrators:

		if chatID == 0 {
			return fmt.Errorf("chatID cannot be zero for the scope %s", scope)
		}

		this.Type = string(scope)
		break

	case
		BotCommandScopeTypeChatMember:

		if chatID == 0 {
			return fmt.Errorf("chatID cannot be zero for the scope %s", scope)
		}

		if userID == 0 {
			return fmt.Errorf("userID cannot be zero for the scope %s", scope)
		}

		this.Type = string(scope)
		break
	}

	return nil
}
