package tgs

const (
	ScopeDefault               ScopeType = "default"
	ScopeAllPrivateChats       ScopeType = "all_private_chats"
	ScopeAllGroupChats         ScopeType = "all_group_chats"
	ScopeAllChatAdministrators ScopeType = "all_chat_administrators"
	ScopeChat                  ScopeType = "chat"                // ScopeChat requires "chat_id"
	ScopeChatAdministrators    ScopeType = "chat_administrators" // ScopeChatAdministrators requires "chat_id"
	ScopeChatMember            ScopeType = "chat_member"         // ScopeChatMember requires "chat_id" & "user_id"
)

type ScopeType string

type Scope struct {
	Type   ScopeType `json:"type" yaml:"type"`
	ChatID int64     `json:"chat_id,omitempty" yaml:"chat_id,omitempty"`
	UserID int64     `json:"user_id,omitempty" yaml:"user_id,omitempty"`
}
