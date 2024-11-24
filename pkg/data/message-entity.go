package data

type MessageEntity struct {
	Type          string  `json:"type" yaml:"type"`
	Offset        int     `json:"offset" yaml:"offset"`
	Length        int     `json:"length" yaml:"length"`
	URL           *string `json:"url" yaml:"url"`
	User          *User   `json:"user" yaml:"user"`
	Language      *string `json:"language" yaml:"language"`
	CustomEmojiID *string `json:"custom_emoji_id" yaml:"custom_emoji_id"`
}
