package tgs

type BotCommandScope struct {
	Type   string `json:"type" yaml:"type"`
	ChatID int64  `json:"chat_id,omitempty" yaml:"chat_id,omitempty"`
	UserID int64  `json:"user_id,omitempty" yaml:"user_id,omitempty"`
}
