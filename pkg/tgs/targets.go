package tgs

const (
	ChatTargetTypePrivate    ChatTargetType = "private"
	ChatTargetTypeGroup      ChatTargetType = "group"
	ChatTargetTypeSuperGroup ChatTargetType = "supergroup"
	ChatTargetTypeChannel    ChatTargetType = "channel"
)

type ChatTargetType string

type Targets struct {
	Users []UserTarget `yaml:"users,omitempty"`
	Chats []ChatTarget `yaml:"chats,omitempty"`
	All   bool         `yaml:"all,omitempty"`
}

func NewTargets() *Targets {
	return &Targets{
		Users: make([]UserTarget, 0),
		Chats: make([]ChatTarget, 0),
	}
}

type UserTarget struct {
	ID int64 `yaml:"id"`
}

type ChatTarget struct {
	Type            ChatTargetType `yaml:"type"`
	ID              int64          `yaml:"id"`
	MessageThreadID int            `yaml:"message_thread_id,omitempty"`
	Users           []UserTarget   `yaml:"users,omitempty"`
}
