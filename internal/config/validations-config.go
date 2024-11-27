package config

type ValidationsConfig struct {
	Users []UserValidation `json:"users,omitempty" yaml:"users,omitempty"`
	Chats []ChatValidation `json:"chats,omitempty" yaml:"chats,omitempty"`
}

func NewValidationsConfig() *ValidationsConfig {
	return &ValidationsConfig{
		Users: make([]UserValidation, 0),
		Chats: make([]ChatValidation, 0),
	}
}

type UserValidation struct {
	ID int64 `json:"id" yaml:"id"`
}

type ChatValidation struct {
	ID   int64  `json:"id" yaml:"id"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}
