package botcommand

type Targets struct {
	Users []UserTarget `json:"users,omitempty" yaml:"users,omitempty"`
	Chats []ChatTarget `json:"chats,omitempty" yaml:"chats,omitempty"`
	All   bool         `json:"all,omitempty" json:"all,omitempty"`
}

func NewTargets() *Targets {
	return &Targets{
		Users: make([]UserTarget, 0),
		Chats: make([]ChatTarget, 0),
	}
}

type UserTarget struct {
	ID int64 `json:"id" yaml:"id"`
}

type ChatTarget struct {
	ID              int64  `json:"id" yaml:"id"`
	Type            string `json:"type,omitempty" yaml:"type,omitempty"`
	MessageThreadID int    `json:"message_thread_id,omitempty" yaml:"message_thread_id,omitempty"`
}
