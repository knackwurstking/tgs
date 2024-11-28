package config

type ValidationTargets struct {
	Users []UserValidationTarget `json:"users,omitempty" yaml:"users,omitempty"`
	Chats []ChatValidationTarget `json:"chats,omitempty" yaml:"chats,omitempty"`
	All   bool                   `json:"all,omitempty" json:"all,omitempty"`
}

func NewValidationsConfig() *ValidationTargets {
	return &ValidationTargets{
		Users: make([]UserValidationTarget, 0),
		Chats: make([]ChatValidationTarget, 0),
	}
}

type UserValidationTarget struct {
	ID int64 `json:"id" yaml:"id"`
}

type ChatValidationTarget struct {
	ID              int64  `json:"id" yaml:"id"`
	Type            string `json:"type,omitempty" yaml:"type,omitempty"`
	MessageThreadID int    `json:"message_thread_id,omitempty" yaml:"message_thread_id,omitempty"`
}
